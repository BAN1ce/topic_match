package message_source

import (
	"context"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"go.uber.org/zap"
)

type StoreSource struct {
	topic       string
	messageChan chan *packet.Message
	storeEvent  broker.MessageStoreEvent
	store       broker.TopicMessageStore
}

func NewStoreSource(topic string, store broker.TopicMessageStore, storeEvent broker.MessageStoreEvent) *StoreSource {
	return &StoreSource{
		topic:       topic,
		storeEvent:  storeEvent,
		store:       store,
		messageChan: make(chan *packet.Message, 10),
	}
}

func (s *StoreSource) NextMessages(ctx context.Context, n int, startMessageID string, include bool) ([]*packet.Message, error) {
	var (
		msg []*packet.Message
		err error
	)

	if startMessageID != "" {
		if msg, err = s.readStoreWriteToWriter(ctx, s.topic, startMessageID, n, include); err != nil {
			return nil, err
		}
		if len(msg) != 0 {
			return msg, nil
		}
	}

	var (
		f            func(...interface{})
		ctx1, cancel = context.WithCancel(ctx)
	)

	logger.Logger.Debug("store source next message start listen store event", zap.String("topic", s.topic), zap.String("startMessageID", startMessageID))

	f = func(i ...interface{}) {

		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*event.StoreEventData)
		if !ok {
			logger.Logger.Error("readStoreWriteToWriter parameter type is not storeEventData", zap.Any("i", i))
			return
		}

		if startMessageID == "" {
			msg, err = s.readStoreWriteToWriter(ctx, s.topic, data.MessageID, n, true)
		} else {
			msg, err = s.readStoreWriteToWriter(ctx, s.topic, startMessageID, n, false)
		}
		cancel()
	}
	s.storeEvent.CreateListenMessageStoreEvent(s.topic, f)
	<-ctx1.Done()
	logger.Logger.Debug("store source next message got from store event", zap.Any("message", msg), zap.Error(err))
	s.storeEvent.DeleteListenMessageStoreEvent(s.topic, f)
	return msg, err
}

func (s *StoreSource) readStoreWriteToWriter(ctx context.Context, topic string, id string, size int, include bool) ([]*packet.Message, error) {
	var (
		message, err = s.store.ReadTopicMessagesByID(ctx, topic, id, size, include)
	)
	if err != nil {
		return nil, err
	}
	logger.Logger.Debug("store help read publish message and write to channel",
		zap.String("store", topic),
		zap.String("id", id),
		zap.Int("size", size),
		zap.Bool("include", include),
		zap.Int("got message size", len(message)))

	return message, nil
}
