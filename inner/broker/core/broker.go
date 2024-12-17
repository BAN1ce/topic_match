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
	"github.com/BAN1ce/skyTree/inner/broker/will"
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
	"github.com/BAN1ce/skyTree/pkg/pool"
	"github.com/eclipse/paho.golang/packets"
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
	cancel                 context.CancelFunc
	server                 *server.Server
	state                  *state.State
	userAuth               middleware.UserAuth
	subTree                broker.SubCenter
	messageStore           *message.Wrapper
	keyStore               store.KeyStore
	clientManager          *client.Clients
	conns                  map[net.Conn]struct{}
	sessionManager         session.Manager
	publishPool            *pool.Publish
	publishRetry           facade.RetrySchedule
	preMiddleware          map[byte][]middleware.PacketMiddleware
	handlers               *Handlers
	mux                    sync.RWMutex
	plugins                *plugin.Plugins
	clientKeepAliveMonitor *monitor.KeepAlive
	willMessageMonitor     *will.MessageMonitor
	retain                 retain.Retain
	clientWg               sync.WaitGroup
	event                  event.Driver
}

func NewBroker(option ...Option) *Broker {
	var (
		b = &Broker{
			publishPool:   pool.NewPublish(),
			preMiddleware: make(map[byte][]middleware.PacketMiddleware),
			conns:         make(map[net.Conn]struct{}, 50000),
		}
	)
	b.initMiddleware()
	b.server = server.NewServer(config.GetConfig().Broker.Listen)

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
	b.ctx, b.cancel = context.WithCancel(ctx)
	go func() {
		if err := b.clientKeepAliveMonitor.Start(b.ctx); err != nil {
			logger.Logger.Error().Err(err).Msg("client keep alive monitor start error")
		}
	}()
	if err := b.server.Start(b.ctx); err != nil {
		return err
	}
	b.acceptConn()
	return nil
}

func (b *Broker) acceptConn() {
	for {
		select {
		case <-b.ctx.Done():
			logger.Logger.Info().Msg("broker closing")
			b.mux.RLock()
			var conns = make([]net.Conn, len(b.conns))
			for c := range b.conns {
				conns = append(conns, c)
			}
			b.mux.RUnlock()

			for _, conn := range conns {
				if conn != nil {
					_ = conn.Close()
				} else {
					logger.Logger.Info().Msg("conn is nil")
				}
			}
			logger.Logger.Info().Msg("broker closed, close all client connection")

			return

		case conn, ok := <-b.server.Conn():
			if !ok {
				logger.Logger.Info().Msg("server closed")
				return
			}

			if conn != nil {
				b.mux.Lock()
				b.conns[conn] = struct{}{}
				b.mux.Unlock()
			}

			newClient := client.NewClient(conn,
				client.WithConfig(
					client.Config{
						WindowSize:       10,
						ReadStoreTimeout: 3 * time.Second,
					}),
				client.WithNotifyClose(b),
				client.WithPlugin(b.plugins),
				client.WithRetain(b.retain),
				client.WithStore(b.messageStore),
				client.WithSubCenter(b.subTree),
				client.WithMessageStoreWrapper(b.messageStore),
				client.WithEvent(b.event),
				client.WithMonitorKeepAlive(b.clientKeepAliveMonitor),
			)

			b.clientWg.Add(1)
			go func(c *client.Client) {
				c.Run(b.ctx, b)
				b.OnClientClose(c)
				b.mux.Lock()
				delete(b.conns, c.GetConn())
				b.mux.Unlock()
				b.clientWg.Done()
				logger.Logger.Info().Str("client", c.MetaString()).Msg("client closed")
			}(newClient)
		}
	}
}

func (b *Broker) Close() error {
	b.cancel()
	err := b.server.Close()

	b.clientWg.Wait()
	return err

}

func (b *Broker) OnClientClose(c *client.Client) {
	b.mux.Lock()
	defer b.mux.Unlock()

	if err := b.clientKeepAliveMonitor.DeleteClient(context.TODO(), c.GetID()); err != nil {
		logger.Logger.Error().Err(err).Str("uid", c.GetID()).Msg("delete client keep alive monitor error")
	}
}

func (b *Broker) OnClose(c *client.Client) {
	logger.Logger.Debug().Str("client", c.MetaString()).Msg("client close")
}

// ------------------------------------ handle client MQTT Packet ------------------------------------//

func (b *Broker) HandlePacket(ctx context.Context, packet *packets.ControlPacket, client *client.Client) (err error) {
	// TODO : check MQTT version
	if err = b.executePreMiddleware(client, packet); err != nil {
		return
	}
	logger.Logger.Debug().Str("packet type", packet.PacketType()).Str("client", client.MetaString()).Msg("handle received packet")

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
		logger.Logger.Warn().Uint8("type", packet.FixedHeader.Type).Msg("unknown packet type")
		err = fmt.Errorf("unknown packet type = %d", packet.FixedHeader.Type)
	}
	return
}

func (b *Broker) executePreMiddleware(client *client.Client, packet *packets.ControlPacket) error {
	for _, midHandle := range b.preMiddleware[packet.FixedHeader.Type] {
		if err := midHandle.Handle(client, packet); err != nil {
			logger.Logger.Error().Err(err).Str("client", client.MetaString()).Any("packet", packet).Msg("middleware handle error")
			client.Close()
			return err
		}
	}
	return nil
}

// ----------------------------------------- support ---------------------------------------------------//
// writePacket for collect all error log
func (b *Broker) writePacket(client *client.Client, packet packets.Packet) {
	if err := client.Write(client2.NewWritePacket(packet)); err != nil {
		logger.Logger.Warn().Err(err).Str("client", client.MetaString()).Msg("write packet error")

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

	b.clientManager.DestroyClient(uid)

	if err := b.clientKeepAliveMonitor.DeleteClient(ctx, uid); err != nil {
		logger.Logger.Error().Err(err).Str("uid", uid).Msg("delete client keep alive monitor error")
	}

}

func (b *Broker) NotifyClientClose(client *client.Client) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.deleteClient(client.GetID())
}

func (b *Broker) NotifyWillMessage(willMessage *session.WillMessage) {
	b.willMessageMonitor.PublishWillMessage(b.ctx, willMessage)

}

func (b *Broker) ReleaseSession(clientID string) {
	var (
		clientSession, ok = b.sessionManager.ReadClientSession(context.TODO(), clientID)
	)
	if !ok {
		return
	}

	if willMessage, ok, err := clientSession.GetWillMessage(); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", clientID).Msg("get will message failed")
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
