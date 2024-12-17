package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/api"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/BAN1ce/skyTree/inner/broker/core"
	"github.com/BAN1ce/skyTree/inner/broker/monitor"
	"github.com/BAN1ce/skyTree/inner/broker/retain"
	"github.com/BAN1ce/skyTree/inner/broker/session"
	"github.com/BAN1ce/skyTree/inner/broker/state"
	"github.com/BAN1ce/skyTree/inner/broker/store"
	"github.com/BAN1ce/skyTree/inner/broker/store/message"
	"github.com/BAN1ce/skyTree/inner/broker/sub_center"
	"github.com/BAN1ce/skyTree/inner/broker/will"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/inner/event/driver"
	"github.com/BAN1ce/skyTree/inner/facade"
	"github.com/BAN1ce/skyTree/inner/version"
	"github.com/BAN1ce/skyTree/logger"
	broker2 "github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/message/serializer"
	"github.com/BAN1ce/skyTree/pkg/broker/plugin"
	"io"
	"log"
	"sync"
	"time"
)

type component interface {
	Start(ctx context.Context) error
	Name() string
}

type App struct {
	ctx        context.Context
	cancel     context.CancelFunc
	brokerCore *core.Broker
	components []component
	mux        sync.Mutex
	wg         sync.WaitGroup
	closer     []io.Closer
}

func NewApp() *App {
	// todo: get config
	//var (
	//	subTree = app2.NewApp()
	//)
	//if err := subTree.StartTopicCluster(context.TODO(), []app2.Option{
	//	app2.WithInitMember(config.GetTree().InitNode),
	//	app2.WithJoin(false),
	//	app2.WithNodeConfig(config.GetTree().NodeHostConfig),
	//	app2.WithConfig(config.GetTree().DragonboatConfig),
	//	app2.WithStateMachine(store2.NewState()),
	//}...); err != nil {
	//	log.Fatal("start store cluster failed", err)
	//}

	var (
		eventDriver = driver.NewLocal()
	)
	event.SetEventDriverOnce(eventDriver)
	var (
		clientManager = client.NewClients()
		sessionStore  = newStore("session")
		messageStore  = newStore("message")
		stateStore    = newStore("state")
		brokerState   = state.NewState(stateStore)

		sessionManager = session.NewSessions(sessionStore)
		storeWrapper   = message.NewStoreWrapper(messageStore, event.StoreEvent)

		subCenter   = broker2.NewLocalSubCenter()
		willMonitor = will.NewWillMessageMonitor(eventDriver, brokerState, subCenter, storeWrapper, sessionManager)
	)

	store.SetDefault(storeWrapper, event.MessageEvent, serializer.ProtoBufVersion)

	// TODO: fix this
	//event.Message.CreateListenClientOffline(shareTopicManger.CloseClient)

	var (
		publishRetrySchedule = facade.InitPublishRetry(clientManager.RetryPublish, clientManager.RetryPublishTimeout)

		brokerCore = core.NewBroker(
			core.WithRetainStore(retain.NewRetainDB(sessionStore)),
			core.WithKeyStore(sessionStore),
			core.WithWillMessageMonitor(willMonitor),
			core.WithPlugins(plugin.NewDefaultPlugin()),
			core.WithPublishRetry(publishRetrySchedule),
			core.WithMessageStore(storeWrapper),
			core.WithState(brokerState),
			core.WithSessionManager(sessionManager),
			core.WithClientManager(clientManager),
			core.WithSubCenter(sub_center.NewSubCenterWrapper(subCenter)),
			core.WithEvent(eventDriver),
			core.WithHandlers(&core.Handlers{
				Connect:     core.NewConnectHandler(),
				Publish:     core.NewPublishHandler(),
				PublishAck:  core.NewPublishAck(),
				PublishRec:  core.NewPublishRec(),
				PublishRel:  core.NewPublishRel(),
				PublishComp: core.NewPublishComp(),
				Ping:        core.NewPingHandler(),
				Sub:         core.NewSubHandler(),
				UnSub:       core.NewUnsubHandler(),
				Auth:        core.NewAuthHandler(),
				Disconnect:  core.NewDisconnectHandler(),
			}))
		app = &App{
			brokerCore: brokerCore,
		}
	)
	app.components = []component{
		app.brokerCore,
		willMonitor,

		monitor.NewMessage(time.Second*30, messageStore),

		api.NewAPI(fmt.Sprintf(":%d", config.GetConfig().Server.Port), &api.Component{
			TopicManager:   brokerCore,
			ClientManager:  clientManager,
			MessageStore:   messageStore,
			SessionManager: sessionManager,
		}),
	}
	app.closer = []io.Closer{
		app.brokerCore,
		messageStore,
		stateStore,
		sessionStore,
	}
	return app
}

func (a *App) Start(ctx context.Context) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.ctx, a.cancel = context.WithCancel(ctx)
	for _, c := range a.components {
		a.wg.Add(1)
		go func(c component) {
			logger.Logger.Info().Str("component", c.Name()).Msg("starting")
			if err := c.Start(ctx); err != nil {
				log.Fatalln("start component", c.Name(), "error: ", err)
			}
			a.wg.Done()
			logger.Logger.Info().Str("component", c.Name()).Msg("closed")
		}(c)
	}
}

func (a *App) Close() error {
	a.mux.Lock()
	defer a.mux.Unlock()
	var err error
	a.cancel()

	for _, c := range a.closer {
		if err1 := c.Close(); err1 != nil {
			err = errors.Join(err, err1)
		}
	}

	logger.Logger.Info().Msg("All closes executed")

	a.wg.Wait()
	return err

}

func (a *App) Version() string {
	return version.Version
}
