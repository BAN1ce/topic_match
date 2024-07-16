package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"go.uber.org/zap"
)

const (
	MessageStored = "event.store.stored"

	MessageRead = "event.store.read"

	MessageDelete = "event.store.delete"
)

var (
	StoreEvent = newStoreEvent()
)

type Store struct {
}

func newStoreEvent() *Store {
	s := &Store{}
	return s
}

func (s *Store) EmitStored(data *StoreEventData) {
	eventDriver.Emit(topicMessageStoredEventName(data.Topic), data)
	eventDriver.Emit(MessageStored, data)

}

func (s *Store) EmitRead(data *StoreEventData) {
	eventDriver.Emit(MessageRead, data)
}

func (s *Store) EmitDelete(data *StoreEventData) {
	eventDriver.Emit(MessageDelete, data)
}

func (s *Store) CreateListenMessageStoreEvent(topic string, handler func(...interface{})) {
	eventDriver.AddListener(topicMessageStoredEventName(topic), handler)
}

func (s *Store) DeleteListenMessageStoreEvent(topic string, handler func(i ...interface{})) {
	eventDriver.RemoveListener(topicMessageStoredEventName(topic), handler)
}

func (s *Store) addDefaultListener() {
	eventDriver.AddListener(MessageRead, func(i ...interface{}) {
		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*StoreEventData)
		if !ok {
			return
		}

		metric.StoreReadDuration.Observe(data.Duration.Seconds())
		metric.StoreReadRequestCount.Inc()

		if ok {
			metric.StoreReadMessageCount.Add(float64(data.Count))
		} else {
			metric.StoreReadRequestFailedCount.Inc()
		}

	})

	eventDriver.AddListener(MessageStored, func(i ...interface{}) {

		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*StoreEventData)
		if !ok {
			return
		}

		metric.StoreWriteDuration.Observe(data.Duration.Seconds())
		metric.StoreWriteRequestCount.Inc()

		if !data.Success {
			metric.StoreWriteRequestFailedCount.Inc()
		}
	})

	eventDriver.AddListener(MessageDelete, func(i ...interface{}) {
		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*StoreEventData)
		if !ok {
			return
		}
		if !data.Success {
			metric.StoreDeleteRequestFailedCount.Inc()
		}
		metric.StoreDeleteDuration.Observe(data.Duration.Seconds())
		metric.StoreDeleteRequestCount.Inc()
	})

}

func topicMessageStoredEventName(topic string) string {
	return utils.WithEventPrefix(MessageStored, topic)
}
