package message_source

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
)

type StoreSource struct {
	topic            string
	store            broker.TopicMessageStore
	publishDoneEvent broker.HandlePublishDoneEvent
}

func NewStoreSource(topic string, store broker.TopicMessageStore, publishDoneEvent broker.HandlePublishDoneEvent) *StoreSource {
	return &StoreSource{
		topic:            topic,
		store:            store,
		publishDoneEvent: publishDoneEvent,
	}
}

func (s *StoreSource) NextMessages(ctx context.Context, n int, startMessageID string, include bool, in chan *packet.Message) (msg []*packet.Message, err error) {
	defer func() {
		for i := 0; i < len(msg); i++ {
			logger.Logger.Debug().Str("start message id", startMessageID).Str("store_message_id", msg[i].MessageID).Str("payload", string(msg[i].PublishPacket.Payload)).Msg("store source next messages")
		}

	}()

	if startMessageID != "" {
		if msg, err = s.readStoreWriteToWriter(ctx, s.topic, startMessageID, n, include); err != nil {
			return nil, err
		}
		if len(msg) != 0 {
			return msg, nil
		}
	}

	var (
		//eventData    *event.StoreEventData
		eventDataCh = make(chan string, 1)

		qos0EventHandler = func(i ...interface{}) {
			if len(i) != 2 {
				logger.Logger.Error().Msg("readStoreWriteToWriter error invalid parameters")
				return
			}

			if msg, ok := i[1].(*packet.Message); ok {
				select {
				case eventDataCh <- msg.MessageID:
				default:

				}
			}
		}
	)

	//s.storeEvent.CreateListenMessageStoreEvent(s.topic, storeEventHandler)
	s.publishDoneEvent.CreateListenPublishEventOnce(s.topic, qos0EventHandler)

	// defer delete listen event
	defer func() {
		//s.storeEvent.DeleteListenMessageStoreEvent(s.topic, storeEventHandler)
		s.publishDoneEvent.DeleteListenPublishEvent(s.topic, qos0EventHandler)
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case m := <-in:
		return []*packet.Message{m}, nil
	case messageID, ok := <-eventDataCh:
		if !ok {
			return nil, ctx.Err()
		}
		if messageID == "" {
			return
		}

		if startMessageID == "" {
			msg, err = s.readStoreWriteToWriter(ctx, s.topic, messageID, n, true)
		} else {
			msg, err = s.readStoreWriteToWriter(ctx, s.topic, startMessageID, n, include)
		}
	}

	return msg, err
}

func (s *StoreSource) readStoreWriteToWriter(ctx context.Context, topic string, id string, size int, include bool) ([]*packet.Message, error) {
	var (
		message, err = s.store.ReadTopicMessagesByID(ctx, topic, id, size, include)
	)
	if err != nil {
		return nil, err
	}

	return message, nil
}
