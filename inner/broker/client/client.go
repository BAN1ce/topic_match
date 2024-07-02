package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/inner/broker/client/sub_topic"
	"github.com/BAN1ce/skyTree/inner/facade"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	topic2 "github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/model"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
	"github.com/BAN1ce/skyTree/pkg/rate"
	"github.com/BAN1ce/skyTree/pkg/retry"
	"github.com/BAN1ce/skyTree/pkg/state"
	"github.com/eclipse/paho.golang/packets"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net"
	"strings"
	"sync"
	"time"
)

type Config struct {
	WindowSize       int
	ReadStoreTimeout time.Duration
	KeepAlive        time.Duration
}

type Handler interface {
	HandlePacket(context.Context, *packets.ControlPacket, *Client) error
}

type Client struct {
	mux  sync.RWMutex
	ctx  context.Context
	conn net.Conn

	ID  string `json:"id"`
	UID string `json:"UID"`

	connectProperties *session.ConnectProperties

	cancel context.CancelCauseFunc

	handler []Handler

	state state.State

	component *component

	packetIDFactory client.PacketIDGenerator

	publishBucket *rate.Bucket
	messages      chan pkg.Message

	topicManager *sub_topic.Topics

	shareTopic map[string]struct{}

	packetIdentifierIDTopic *PacketIDTopic

	QoS2 *QoS2ReceiveStore

	topicAliasFromClient map[uint16]string

	willFlag  bool
	keepAlive time.Duration
	aliveTime time.Time
}

func NewClient(conn net.Conn, option ...Component) *Client {
	var (
		c = &Client{
			UID:                     uuid.NewString(),
			conn:                    conn,
			component:               new(component),
			packetIdentifierIDTopic: NewPacketIDTopic(),
			topicAliasFromClient:    map[uint16]string{},
			shareTopic:              map[string]struct{}{},
		}
	)
	for _, o := range option {
		o(c.component)
	}
	c.packetIDFactory = NewPacketIDFactory()
	c.messages = make(chan pkg.Message, c.component.cfg.WindowSize)
	c.QoS2 = NewQoS2Handler()

	return c
}

func (c *Client) Run(ctx context.Context, handler Handler) {
	ctx = context.WithValue(ctx, pkg.ClientUIDKey, c.UID)
	c.ctx, c.cancel = context.WithCancelCause(ctx)
	c.topicManager = sub_topic.NewTopics(c.ctx)
	c.handler = []Handler{
		handler,
		newClientHandler(c),
	}

	var (
		controlPacket *packets.ControlPacket
		err           error
	)

	for {
		controlPacket, err = packets.ReadPacket(c.conn)
		if err != nil {
			logger.Logger.Info("read controlPacket error = ", zap.Error(err), zap.String("client", c.MetaString()))
			c.afterClose()
			return
		}
		for _, handler := range c.handler {
			if err := handler.HandlePacket(context.TODO(), controlPacket, c); err != nil {
				logger.Logger.Warn("handle packet error", zap.Error(err), zap.String("client", c.MetaString()))
				c.CloseIgnoreError()
				break
			}
		}
	}
}

func (c *Client) AddTopics(topic map[string]topic2.Topic) {
	c.mux.Lock()
	defer c.mux.Unlock()
	for topicName, t := range topic {
		if err := c.topicManager.AddTopic(topicName, t); err != nil {
			logger.Logger.Error("add topic error", zap.Error(err), zap.String("client", c.MetaString()), zap.Any("topic", t))
		}
	}
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if errors.Is(err, net.ErrClosed) {
		return nil
	}
	return err
}

func (c *Client) CloseIgnoreError() {
	if err := c.conn.Close(); err != nil {
		logger.Logger.Warn("close conn error", zap.Error(err), zap.String("client", c.MetaString()))
	}
}

func (c *Client) afterClose() {
	logger.Logger.Info("client close", zap.String("clientID", c.ID))
	defer func() {
		c.component.notifyClose.NotifyClientClose(c)
	}()

	if err := c.conn.Close(); err != nil {
		logger.Logger.Info("close conn error", zap.Error(err), zap.String("client", c.MetaString()))
	}

	//  close normal topicManager
	if err := c.topicManager.Close(); err != nil {
		logger.Logger.Warn("close topicManager error", zap.Error(err), zap.String("client", c.MetaString()))
	}

	//  close share topicManager

	if c.component.session == nil {
		logger.Logger.Warn("session is nil", zap.String("client", c.MetaString()))
		return
	}

	// store unfinished message to session for qos1 and qos2
	c.storeUnfinishedMessage()

	//  handle will message
	willMessage, ok, err := c.component.session.GetWillMessage()
	if !ok {
		return
	}
	if err != nil {
		logger.Logger.Error("get will message error", zap.Error(err), zap.String("client", c.MetaString()))
		return
	}

	var (
		willDelayInterval = willMessage.Property.GetDelayInterval()
	)

	if willDelayInterval == 0 {
		c.component.notifyClose.NotifyWillMessage(willMessage)
		return
	}

	// create a delay task to publish will message
	if err := facade.GetWillDelay().Create(&retry.Task{
		MaxTimes:     1,
		MaxTime:      0,
		IntervalTime: willDelayInterval,
		Key:          willMessage.DelayTaskID,
		Data:         willMessage,
		Job: func(task *retry.Task) {
			if m, ok := task.Data.(*session.WillMessage); ok {
				logger.Logger.Debug("will message delay task", zap.String("client", c.MetaString()), zap.String("delayTaskID", task.Key))
				c.component.notifyClose.NotifyWillMessage(m)
			} else {
				logger.Logger.Error("will message type error", zap.String("client", c.MetaString()))
			}
		},
		TimeoutJob: nil,
	}); err != nil {
		logger.Logger.Error("create will delay task error", zap.Error(err), zap.String("client", c.MetaString()))
		return
	}

}

func (c *Client) storeUnfinishedMessage() {
	var message = c.topicManager.GetUnfinishedMessage()
	if c.component.session == nil {
		return
	}
	for topicName, message := range message {
		logger.Logger.Debug("store unfinished message", zap.String("client", c.MetaString()), zap.String("topicName", topicName), zap.Int("message", len(message)))
		c.component.session.CreateTopicUnFinishedMessage(topicName, message)
	}
}

func (c *Client) SetID(id string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.ID = id
}

func (c *Client) WritePacket(packet *client.WritePacket) error {
	if packet == nil {
		return fmt.Errorf("packet is nil")
	}

	if _, ok := packet.Packet.(*packets.Publish); ok {
		//  get token before a written packet
		c.publishBucket.GetToken()
	}

	// try lock after get token
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.writePacket(packet)
}

func (c *Client) writePacket(packet *client.WritePacket) error {
	var (
		buf              = pool.ByteBufferPool.Get()
		prepareWriteSize int64
		err              error
		topicName        string
	)

	defer pool.ByteBufferPool.Put(buf)
	// publishAck, subscribeAck, unsubscribeAck should use the same packetID as the original packet

	switch p := packet.Packet.(type) {

	case *packets.Connack:
		// do plugin
		c.component.plugin.DoSendConnAck(c.ID, p)

	case *packets.Suback:
		// do plugin
		c.component.plugin.DoSendSubAck(c.ID, p)

	case *packets.Unsuback:
		// do plugin
		c.component.plugin.DoSendUnsubAck(c.ID, p)

	case *packets.Publish:
		// do plugin
		c.component.plugin.DoSendPublish(c.ID, p)
		c.packetIdentifierIDTopic.SePacketIDTopic(p.PacketID, packet.FullTopic)
		logger.Logger.Debug("publish to client", zap.Any("packet", p), zap.Uint16("packetID", p.PacketID), zap.String("client", c.MetaString()), zap.String("store", p.Topic))

	case *packets.Puback:
		// do plugin
		c.component.plugin.DoSendPubAck(c.ID, p)

	case *packets.Pubrec:
		// do plugin
		logger.Logger.Debug("send pubrec", zap.String("client", c.MetaString()), zap.Uint16("packetID", p.PacketID))
		c.component.plugin.DoSendPubRec(c.ID, p)

	case *packets.Pubrel:
		// do plugin
		c.component.plugin.DoSendPubRel(c.ID, p)

	case *packets.Pubcomp:
		// do plugin
		c.component.plugin.DoSendPubComp(c.ID, p)

	case *packets.Pingresp:
		// do plugin
		c.component.plugin.DoSendPingResp(c.ID, p)

	}

	if prepareWriteSize, err = packet.Packet.WriteTo(buf); err != nil {
		logger.Logger.Info("write packet error", zap.Error(err), zap.String("client", c.MetaString()))
		// TODO: check maximum packet size should close client ?
	}

	// check maximum packet size
	if c.connectProperties != nil && c.connectProperties.MaximumPacketSize != nil && *c.connectProperties.MaximumPacketSize != 0 {
		if buf.Len() > int(*c.connectProperties.MaximumPacketSize) && *c.connectProperties.MaximumPacketSize != 0 {
			return errs.ErrPacketOversize
		}
	}

	if _, err = c.conn.Write(buf.Bytes()); err != nil {
		logger.Logger.Info("write packet error", zap.Error(err), zap.String("client", c.MetaString()), zap.Int64("prepareWriteSize", prepareWriteSize), zap.String("topicName", topicName))
		// TODO: check maximum packet size should close client ?
		if err := c.conn.Close(); err != nil {
			logger.Logger.Warn("close conn error", zap.Error(err), zap.String("client", c.MetaString()))
		}
		return err
	}
	return nil
}

func (c *Client) SetComponent(ops ...Component) error {
	c.mux.Lock()
	for _, o := range ops {
		o(c.component)
	}

	c.mux.Unlock()
	return nil
}

func (c *Client) getSession() session.Session {
	return c.component.session
}

func (c *Client) setWill(message *session.WillMessage) error {
	c.willFlag = true
	if oldWillMessage, ok, _ := c.component.session.GetWillMessage(); ok {
		facade.GetWillDelay().Delete(oldWillMessage.DelayTaskID)
	}
	return c.component.session.SetWillMessage(message)
}
func (c *Client) DeleteWill() error {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.component.session.DeleteWillMessage()
}

func (c *Client) SetClientTopicAlias(topic string, alias uint16) {
	c.mux.Lock()
	c.topicAliasFromClient[alias] = topic
	c.mux.Unlock()
}

func (c *Client) SetConnectProperties(properties *session.ConnectProperties) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.connectProperties = properties
	return c.component.session.SetConnectProperties(properties)
}
func (c *Client) GetConnectProperties() session.ConnectProperties {
	c.mux.RLock()
	defer c.mux.RUnlock()
	if c.connectProperties == nil {
		return session.ConnectProperties{}
	}
	return *c.connectProperties
}

func (c *Client) GetID() string {
	return c.ID
}

func (c *Client) GetUid() string {
	return c.UID
}

func (c *Client) MetaString() string {
	var (
		s strings.Builder
	)
	s.WriteString("clientID: ")
	s.WriteString(c.ID)
	s.WriteString(" ")
	s.WriteString("remoteAddr: ")
	s.WriteString(c.conn.RemoteAddr().String())
	return s.String()
}

func (c *Client) Meta() *model.Meta {
	return &model.Meta{
		ID:        c.ID,
		AliveTime: c.aliveTime,
		KeepAlive: c.keepAlive.String(),
		Topic:     c.topicManager.Meta(),
	}
}

func (c *Client) GetClientTopicAlias(u uint16) string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.topicAliasFromClient[u]
}

func (c *Client) Publish(topic string, message *packet.Message) error {
	return c.topicManager.Publish(topic, message)
}

func (c *Client) RefreshAliveTime() {
	c.aliveTime = time.Now()
}
func (c *Client) GetKeepAliveTime() time.Duration {
	return c.keepAlive
}

func (c *Client) GetSession() session.Session {
	return c.component.session
}

func (c *Client) NextPacketID() uint16 {
	return c.packetIDFactory.NextPacketID()
}
