package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/inner/broker/client/sub_topic"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/inner/facade"
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/model"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
	"github.com/BAN1ce/skyTree/pkg/rate"
	"github.com/BAN1ce/skyTree/pkg/retry"
	"github.com/BAN1ce/skyTree/pkg/state"
	"github.com/eclipse/paho.golang/packets"
	"github.com/zyedidia/generic/list"
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

	ID string `json:"id"`

	Username string `json:"username"`

	connectProperties *session.ConnectProperties

	cancel context.CancelCauseFunc

	handler []Handler

	state state.State

	component *Component

	packetIDFactory client.PacketIDGenerator

	publishBucket *rate.Bucket

	messages chan pkg.Message

	topicManager *sub_topic.Topics

	shareTopic map[string]struct{}

	packetIdentifierIDTopic *PacketIDTopic

	QoS2 *QoS2ReceiveStore

	topicAliasFromClient map[uint16]string

	willFlag  bool
	keepAlive time.Duration
	aliveTime time.Time

	receiveQos2 *list.List[*packet.Message]
}

func NewClient(conn net.Conn, option ...ComponentOption) *Client {
	var (
		c = &Client{
			conn:                    conn,
			component:               new(Component),
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
	c.QoS2 = NewQoS2ReceiveStore()

	return c
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

func (c *Client) Run(ctx context.Context, handler Handler) {
	ctx = context.WithValue(ctx, pkg.ClientIDKey, c.GetID())
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
	defer func() {
		c.afterClose()
	}()

	for {
		select {
		// maybe never trigger, because the client blocked in read
		case <-c.ctx.Done():
			return
		default:

			controlPacket, err = packets.ReadPacket(c.conn)
			if err != nil {
				logger.Logger.Info().Str("client", c.MetaString()).Err(err).Msg("read controlPacket error")
				return
			}

			for _, handler := range c.handler {
				if err := handler.HandlePacket(context.TODO(), controlPacket, c); err != nil {
					logger.Logger.Warn().Str("client", c.MetaString()).Err(err).Msg("handle packet error")
					c.CloseIgnoreError()
					break
				}
			}
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
		logger.Logger.Info().Str("client", c.MetaString()).Err(err).Msg("close conn error")
	}
}

func (c *Client) afterClose() {
	defer func() {
		c.component.notifyClose.NotifyClientClose(c)
	}()

	if err := c.conn.Close(); err != nil {
		logger.Logger.Warn().Str("client", c.MetaString()).Err(err).Msg("close conn error")
	}

	//  close normal topicManager
	if err := c.topicManager.Close(); err != nil {
		logger.Logger.Warn().Err(err).Str("client", c.MetaString()).Msg("close topicManager error")
	}

	// store all unfinished messages to session for qos1 and qos2
	c.storeUnfinishedMessage()

	if c.component.session == nil {
		logger.Logger.Warn().Str("client", c.MetaString()).Msg("session is nil")
		return
	}

	//  handle will message
	willMessage, ok, err := c.component.session.GetWillMessage()

	if !ok {
		return
	}

	if err != nil {
		logger.Logger.Error().Err(err).Str("client", c.MetaString()).Msg("get will message error")
		return
	}

	c.component.notifyClose.NotifyWillMessage(willMessage)

}

func (c *Client) storeUnfinishedMessage() {
	var message = c.topicManager.GetUnfinishedMessage()
	if c.component.session == nil {
		return
	}
	for topicName, message := range message {
		c.component.session.CreateTopicUnFinishedMessage(topicName, message)
	}
}

func (c *Client) SetID(id string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.ID = id
}

func (c *Client) Write(packet *client.WritePacket) error {
	if packet == nil {
		return fmt.Errorf("packet is nil")
	}

	if tmp, ok := packet.Packet.(*packets.Publish); ok {
		//  get token before a written packet
		if tmp.QoS != broker.QoS0 {
			c.publishBucket.GetToken(c.ctx)
		}

		logger.Logger.Debug().Str("client", c.ID).Any("publish packet", packet.Packet).Str("full topic", packet.FullTopic).Uint16("packetID", tmp.PacketID).Msg("write packet")
	}

	// try lock after get token
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.writePacket(packet)
}

func (c *Client) RetryWrite(packet *client.WritePacket) error {
	if packet == nil {
		return fmt.Errorf("packet is nil")
	}

	// try lock after get token

	if p, ok := packet.Packet.(*packets.Publish); ok {
		pub := pool.PublishPool.Get()
		defer pool.PublishPool.Put(pub)
		pool.CopyPublish(pub, p)
		pub.Duplicate = true
		packet.Packet = pub
	}
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
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.CONNACK)
		// do plugin
		c.component.plugin.DoSendConnAck(c.ID, p)

	case *packets.Suback:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.SUBACK)
		// do plugin
		c.component.plugin.DoSendSubAck(c.ID, p)

	case *packets.Unsuback:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.UNSUBACK)
		// do plugin
		c.component.plugin.DoSendUnsubAck(c.ID, p)

	case *packets.Publish:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.PUBLISH)
		// do plugin
		c.component.plugin.DoSendPublish(c.ID, p)
		c.packetIdentifierIDTopic.SePacketIDTopic(p.PacketID, packet.FullTopic)

	case *packets.Puback:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.PUBACK)
		// do plugin
		c.component.plugin.DoSendPubAck(c.ID, p)

	case *packets.Pubrec:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.PUBREC)
		// do plugin
		c.component.plugin.DoSendPubRec(c.ID, p)

	case *packets.Pubrel:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.PUBREL)
		// do plugin
		c.component.plugin.DoSendPubRel(c.ID, p)

	case *packets.Pubcomp:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.PUBCOMP)
		// do plugin
		c.component.plugin.DoSendPubComp(c.ID, p)

	case *packets.Pingresp:
		event.MQTTEvent.EmitSendMQTTPacketEvent(packets.PINGRESP)
		// do plugin
		c.component.plugin.DoSendPingResp(c.ID, p)

	}

	if prepareWriteSize, err = packet.Packet.WriteTo(buf); err != nil {
		logger.Logger.Info().Err(err).Str("client", c.MetaString()).Int64("write size", prepareWriteSize).Str("topic", topicName).Msg("write packet error")
		// TODO: check maximum packet size should close client ?
	}

	// check maximum packet size
	if c.connectProperties != nil && c.connectProperties.MaximumPacketSize != nil && *c.connectProperties.MaximumPacketSize != 0 {
		if buf.Len() > int(*c.connectProperties.MaximumPacketSize) && *c.connectProperties.MaximumPacketSize != 0 {
			return errs.ErrPacketOversize
		}
	}

	if _, err = c.conn.Write(buf.Bytes()); err != nil {
		logger.Logger.Info().Err(err).Str("client", c.MetaString()).Msg("write packet error")
		// TODO: check maximum packet size should close client ?
		if err := c.conn.Close(); err != nil {
			logger.Logger.Warn().Err(err).Str("client", c.MetaString()).Msg("close conn error")
		}
		return err
	}
	return nil
}

func (c *Client) SetComponent(ops ...ComponentOption) error {
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

func (c *Client) GetClientTopicAlias(u uint16) string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	return c.topicAliasFromClient[u]
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

func (c *Client) RetrySend(message *packet.Message) error {
	var topic = message.PublishPacket.Topic
	defer func() {
		metric.ExecuteRetryTask.Inc()
	}()
	if message.ShareTopic != "" {
		topic = message.ShareTopic
	}
	if !message.RetryInfo.IsTimeout() {
		message.RetryInfo.Times.Add(1)
	}
	logger.Logger.Info().Msg("Client retry publish")
	// find the topic instance, and resend

	//  create a retry task
	if err := facade.GetPublishRetry().Create(retry.NewTask(message.RetryInfo.Key, message, c.ID, message.RetryInfo.IntervalTime)); err != nil {
		logger.Logger.Error().Err(err).Msg("WithRetryClient: publish retry failed")
	}

	if message.PubReceived {
		pubRel := pool.NewPubRel().Get()
		defer pool.NewPubRel().Put(pubRel)
		pubRel.PacketID = message.PublishPacket.PacketID
		pubRel.Properties = message.PublishPacket.Properties
		return c.RetryWrite(&client.WritePacket{
			Packet: pubRel,
		})
	}

	return c.RetryWrite(&client.WritePacket{
		Packet:    message.PublishPacket,
		FullTopic: topic,
	})
	//logger.Logger.Debug().Str("client", c.MetaString()).Str("topic", topic).Uint16("packetID", message.PublishPacket.PacketID).Msg("retry publish")
	//return c.topicManager.Publish(topic, message)
}
