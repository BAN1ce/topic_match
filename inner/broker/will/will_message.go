package will

import (
	"context"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/broker/store"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"time"
)

type MessageStore interface {
	GetWillClientID(ctx context.Context, time time.Time) (clientID []string)
	DeleteWillClientID(ctx context.Context, key []string)
	AddDelayWillTask(ctx context.Context, clientID string, triggerTime time.Time) error
}

type MessageMonitor struct {
	willMessageStore    MessageStore
	subCenter           broker.SubCenter
	messageStoreWrapper store.Wrapper
	sessions            session.Manager
	event               event.Driver
}

func NewWillMessageMonitor(event event.Driver, messageStore MessageStore, center broker.SubCenter, wrapper store.Wrapper, manager session.Manager) *MessageMonitor {
	return &MessageMonitor{
		willMessageStore:    messageStore,
		subCenter:           center,
		messageStoreWrapper: wrapper,
		sessions:            manager,
		event:               event,
	}
}

func (w *MessageMonitor) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Logger.Info().Msg("will message monitor exit by context done")
			return
		case <-ticker.C:
			w.ScanDelayAndPublish()
		}
	}
}

func (w *MessageMonitor) ScanDelayAndPublish() {
	var (
		ctx, cancel = context.WithTimeout(context.TODO(), 3*time.Second)
		triggerTime = time.Now()
	)
	defer cancel()

	// FIXME: add buffer to avoid blocking

	clientID := w.willMessageStore.GetWillClientID(ctx, triggerTime)

	for _, id := range clientID {
		w.TriggerClientWill(ctx, id)
	}

	w.willMessageStore.DeleteWillClientID(ctx, clientID)
}

func (w *MessageMonitor) TriggerClientWill(ctx context.Context, clientID string) {
	var (
		ok   bool
		se   session.Session
		will *session.WillMessage
		err  error
	)

	logger.Logger.Debug().Str("clientID", clientID).Msg("trigger will message")

	se, ok = w.sessions.ReadClientSession(ctx, clientID)
	if !ok {
		logger.Logger.Warn().Str("clientID", clientID).Msg("client session not found for will message")
		return
	}

	will, ok, err = se.GetWillMessage()
	if err != nil || !ok {
		logger.Logger.Error().Err(err).Str("clientID", clientID).Msg("get will message error")
		return
	}

	w.publishWill(will)

}

func (w *MessageMonitor) PublishWillMessage(ctx context.Context, willMessage *session.WillMessage) {
	logger.Logger.Debug().Str("clientID", willMessage.ClientID).Msg("try to publish will message")

	// if property is nil or willDelayInterval is nil, we publish they will message immediately
	if willMessage.Property == nil || willMessage.Property.WillDelayInterval == nil {
		w.publishWill(willMessage)
	}

	// if willDelayInterval is 0, we publish they will message immediately
	if *willMessage.Property.WillDelayInterval == 0 {
		w.publishWill(willMessage)
	}

	// if willDelayInterval is not 0, we add a delay task
	if err := w.willMessageStore.AddDelayWillTask(ctx, willMessage.ClientID, time.Now().Add(time.Duration(*willMessage.Property.WillDelayInterval)*time.Second)); err != nil {
		logger.Logger.Error().Err(err).Msg("add delay will task error")
	}

}

func (w *MessageMonitor) publishWill(will *session.WillMessage) {
	logger.Logger.Debug().Str("clientID", will.ClientID).Msg("publish will message")

	var (
		message = &packet.Message{
			PublishPacket: will.ToPublishPacket(),
			ExpiredTime:   0,
			Will:          true,
			State:         0,
			OwnerClientID: will.ClientID,
		}
	)

	topics := w.subCenter.GetAllMatchTopics(message.PublishPacket.Topic)

	for topic := range topics {
		// for QoS0 message, we don't need to wait for the ack
		event.MessageEvent.EmitFinishTopicPublishMessage(w.event, topic, &packet.Message{
			PublishPacket: will.ToPublishPacket(),
			Will:          true,
			OwnerClientID: will.ClientID,
		})
	}

	if _, err := w.messageStoreWrapper.StorePublishPacket(topics, message); err != nil {
		logger.Logger.Error().Err(err).Msg("willMessageStore will message TriggerClientWill packet error")
	}

}

func (w *MessageMonitor) Start(ctx context.Context) error {
	go w.Run(ctx)
	return nil
}

func (w *MessageMonitor) Name() string {
	return "willMessageMonitor"
}

func (w *MessageMonitor) Close() error {
	return nil
}
