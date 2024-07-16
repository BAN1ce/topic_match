package core

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/BAN1ce/skyTree/inner/broker/monitor"
	"github.com/BAN1ce/skyTree/inner/broker/server"
	"github.com/BAN1ce/skyTree/inner/broker/state"
	"github.com/BAN1ce/skyTree/inner/broker/store/message"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/inner/facade"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	client2 "github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/plugin"
	"github.com/BAN1ce/skyTree/pkg/broker/retain"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/broker/store"
	"github.com/BAN1ce/skyTree/pkg/middleware"
	"github.com/BAN1ce/skyTree/pkg/model"
	packet2 "github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
	"github.com/eclipse/paho.golang/packets"
	"go.uber.org/zap"
	"log"
	"net"
	"sync"
	"time"
)

type Observer interface {
	OnClientClose(b Broker, c *client.Client)
}

type Handlers struct {
	Connect     brokerHandler
	Publish     brokerHandler
	PublishAck  brokerHandler
	PublishRec  brokerHandler
	PublishRel  brokerHandler
	PublishComp brokerHandler
	Ping        brokerHandler
	Sub         brokerHandler
	UnSub       brokerHandler
	Auth        brokerHandler
	Disconnect  brokerHandler
}

type brokerHandler interface {
	Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error)
}

type Broker struct {
	ctx                    context.Context
	server                 *server.Server
	state                  *state.State
	userAuth               middleware.UserAuth
	subTree                broker.SubCenter
	messageStore           *message.Wrapper
	keyStore               store.KeyStore
	clientManager          *client.Clients
	sessionManager         session.Manager
	publishPool            *pool.Publish
	publishRetry           facade.RetrySchedule
	preMiddleware          map[byte][]middleware.PacketMiddleware
	handlers               *Handlers
	mux                    sync.Mutex
	plugins                *plugin.Plugins
	clientKeepAliveMonitor *monitor.KeepAlive
	retain                 retain.Retain
}

func NewBroker(option ...Option) *Broker {
	var (
		ip   = `0.0.0.0`
		port = config.GetServer().GetBrokerPort()
		b    = &Broker{
			publishPool:   pool.NewPublish(),
			preMiddleware: make(map[byte][]middleware.PacketMiddleware),
		}
	)
	b.initMiddleware()
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		log.Fatalln("resolve tcp addr error: ", err)
	}
	b.server = server.NewTCPServer(tcpAddr)

	for _, opt := range option {
		opt(b)
	}

	b.clientKeepAliveMonitor = monitor.NewKeepAlive(b.keyStore, 30*time.Second, b.CloseExpiredClient)
	return b
}

func (b *Broker) Name() string {
	return "messageStore"
}

func (b *Broker) Start(ctx context.Context) error {
	b.ctx = ctx
	go func() {
		if err := b.clientKeepAliveMonitor.Start(b.ctx); err != nil {
			logger.Logger.Error("client keep alive monitor start error", zap.Error(err))
		}
	}()
	if err := b.server.Start(); err != nil {
		return err
	}
	b.acceptConn()
	return nil
}

func (b *Broker) acceptConn() {
	var wg sync.WaitGroup
	for {
		conn, ok := b.server.Conn()
		if !ok {
			logger.Logger.Info("server closed")
			return
		}
		newClient := client.NewClient(conn, client.WithConfig(client.Config{
			WindowSize:       10,
			ReadStoreTimeout: 3 * time.Second,
		}), client.WithNotifyClose(b), client.WithPlugin(b.plugins), client.WithRetain(b.retain), client.WithStore(b.messageStore))

		wg.Add(1)
		go func(c *client.Client) {
			c.Run(b.ctx, b)
			b.OnClientClose(c)
			wg.Done()
			logger.Logger.Info("client closed", zap.String("client", c.MetaString()))
		}(newClient)
	}
}

func (b *Broker) OnClientClose(c *client.Client) {
	b.mux.Lock()
	defer b.mux.Unlock()

	if err := b.clientKeepAliveMonitor.DeleteClient(context.TODO(), c.GetUid()); err != nil {
		logger.Logger.Error("delete client keep alive monitor error", zap.Error(err), zap.String("uid", c.GetUid()))
	}
}

func (b *Broker) OnClose(c *client.Client) {
	logger.Logger.Debug("broker receive client close event", zap.String("client", c.MetaString()))
}

// ------------------------------------ handle client MQTT PublishPacket ------------------------------------//

func (b *Broker) HandlePacket(ctx context.Context, packet *packets.ControlPacket, client *client.Client) (err error) {
	// TODO : check MQTT version
	if err = b.executePreMiddleware(client, packet); err != nil {
		return
	}
	logger.Logger.Debug("Handle Receive Packet", zap.String("client", client.MetaString()), zap.Any("packet", packet))

	// TODO: emit event with packet.Content pointer to avoid copy ?  but need to check if it is safe.
	event.MQTTEvent.EmitReceivedMQTTPacketEvent(packet.FixedHeader.Type)
	switch packet.FixedHeader.Type {
	case packets.CONNECT:
		err = b.handlers.Connect.Handle(b, client, packet)
	case packets.PUBLISH:
		err = b.handlers.Publish.Handle(b, client, packet)

	case packets.PUBACK:
		err = b.handlers.PublishAck.Handle(b, client, packet)
	case packets.PUBREC:
		err = b.handlers.PublishRec.Handle(b, client, packet)
	case packets.PUBREL:
		err = b.handlers.PublishRel.Handle(b, client, packet)
	case packets.PUBCOMP:
		err = b.handlers.PublishComp.Handle(b, client, packet)

	case packets.SUBSCRIBE:
		err = b.handlers.Sub.Handle(b, client, packet)
	case packets.UNSUBSCRIBE:
		err = b.handlers.UnSub.Handle(b, client, packet)
	case packets.PINGREQ:
		err = b.handlers.Ping.Handle(b, client, packet)
	case packets.DISCONNECT:
		err = b.handlers.Disconnect.Handle(b, client, packet)
	case packets.AUTH:
		err = b.handlers.Auth.Handle(b, client, packet)
	default:
		logger.Logger.Warn("unknown packet type = ", zap.Uint8("type", packet.FixedHeader.Type))
		err = fmt.Errorf("unknown packet type = %d", packet.FixedHeader.Type)
	}
	return
}

func (b *Broker) executePreMiddleware(client *client.Client, packet *packets.ControlPacket) error {
	for _, midHandle := range b.preMiddleware[packet.FixedHeader.Type] {
		if err := midHandle.Handle(client, packet); err != nil {
			logger.Logger.Error("middleware handle error: ", zap.Error(err), zap.String("client", client.MetaString()))
			client.Close()
			return err
		}
	}
	return nil
}

// ----------------------------------------- support ---------------------------------------------------//
// writePacket for collect all error log
func (b *Broker) writePacket(client *client.Client, packet packets.Packet) {
	if err := client.WritePacket(client2.NewWritePacket(packet)); err != nil {
		logger.Logger.Warn("write packet error", zap.Error(err), zap.String("client", client.MetaString()))

	}
}

func (b *Broker) CreateClient(client *client.Client) {
	b.mux.Lock()
	b.clientManager.CreateClient(client)

	event.ClientEvent.EmitClientCountState(b.clientManager.Count())
	b.mux.Unlock()

}
func (b *Broker) DeleteClient(uid string) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.deleteClient(uid)
}

func (b *Broker) deleteClient(uid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	event.ClientEvent.EmitClientCountState(b.clientManager.Count())

	event.ClientEvent.EmitClientDeleteSuccess(uid)

	logger.Logger.Debug("broker delete deleteClient", zap.String("uid", uid))

	b.clientManager.DestroyClient(uid)

	if err := b.clientKeepAliveMonitor.DeleteClient(ctx, uid); err != nil {
		logger.Logger.Error("delete deleteClient keep alive monitor error", zap.Error(err), zap.String("uid", uid))
	}

}

func (b *Broker) NotifyClientClose(client *client.Client) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.deleteClient(client.GetUid())
}

func (b *Broker) NotifyWillMessage(willMessage *session.WillMessage) {
	var (
		publishMessage = &packet2.Message{
			PublishPacket: willMessage.ToPublishPacket(),
			PubRelPacket:  nil,
			ExpiredTime:   0,
			Will:          true,
			State:         0,
		}
	)
	logger.Logger.Info("notify will message", zap.String("topic", willMessage.Topic))

	event.MessageEvent.EmitReceivedPublishDone(willMessage.Topic, publishMessage)

	topics := b.subTree.MatchTopic(willMessage.Topic)

	if len(topics) == 0 {
		logger.Logger.Info("notify will message, no client subscribe topic", zap.String("topic", willMessage.Topic))
		return
	}
	if _, err := b.messageStore.StorePublishPacket(topics, publishMessage); err != nil {
		logger.Logger.Error("messageStore will message publish packet error", zap.Error(err), zap.String("topic", willMessage.Topic))
		return
	}

}

func (b *Broker) ReadTopicRetainMessage(topic string) []*packet2.Message {
	var (
		retainMessage  []*packet2.Message
		messageID, err = b.state.ReadRetainMessageID(topic)
		ctx, cancel    = context.WithCancel(b.ctx)
	)
	cancel()
	if err != nil {
		logger.Logger.Error("read retain message error", zap.Error(err), zap.String("topic", topic))
		return nil
	}
	for _, id := range messageID {
		if err := b.messageStore.ReadPublishMessage(ctx, topic, id, 1, true, func(message *packet2.Message) {
			retainMessage = append(retainMessage, message)
		}); err != nil {
			logger.Logger.Error("read retain message error", zap.Error(err), zap.String("topic", topic), zap.String("messageID", id))
		}
	}
	return retainMessage

}
func (b *Broker) ReleaseSession(clientID string) {
	var (
		clientSession, ok = b.sessionManager.ReadClientSession(context.TODO(), clientID)
	)
	if !ok {
		return
	}

	if willMessage, ok, err := clientSession.GetWillMessage(); err != nil {
		logger.Logger.Error("get will message failed", zap.Error(err), zap.String("clientID", clientID))
	} else if ok {
		// delete will message from delay queue
		facade.GetWillDelay().Delete(willMessage.DelayTaskID)
	}
	clientSession.Release()
}

func (b *Broker) ReadTopic(topic string) (*model.Topic, error) {
	var (
		result = &model.Topic{}
	)
	if m, ok := b.retain.GetRetainMessage(topic); ok {
		result.Retain = m
	}
	c := b.subTree.Match(topic)
	for clientID, qos := range c {
		result.SubClient = append(result.SubClient, &model.Subscriber{
			ID:  clientID,
			QoS: int(qos),
		})
	}
	return result, nil
}

// CloseExpiredClient close expired client
func (b *Broker) CloseExpiredClient(uid []string) {
	//logger.Logger.Info("close expired client", zap.Strings("uid", uid))
	//for _, id := range uid {
	//	b.DeleteClient(id)
	//}
}
