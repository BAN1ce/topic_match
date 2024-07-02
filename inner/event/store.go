package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/kataras/go-events"
	"go.uber.org/zap"
	"time"
)

type StoreEventData struct {
	Topic     string
	MessageID string
	Success   bool
	Duration  time.Duration
	Count     int
}

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
	s.registerMetric()
	return s
}

func (s *Store) EmitStored(data *StoreEventData) {
	Driver.Emit(TopicMessageStoredEventName(data.Topic), data)
	Driver.Emit(MessageStored, data)

}

func (s *Store) EmitRead(data *StoreEventData) {
	Driver.Emit(MessageRead, data)
}

func (s *Store) EmitDelete(data *StoreEventData) {
	Driver.Emit(MessageDelete, data)
}

func (s *Store) CreateListenMessageStoreEvent(topic string, handler func(...interface{})) {
	Driver.AddListener(TopicMessageStoredEventName(topic), handler)
}

func (s *Store) DeleteListenMessageStoreEvent(topic string, handler func(i ...interface{})) {
	Driver.RemoveListener(TopicMessageStoredEventName(topic), handler)
}

func (s *Store) registerMetric() {
	Driver.AddListener(MessageRead, func(i ...interface{}) {
		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*StoreEventData)
		if !ok {
			return
		}

		metric.StoreReadDuration.Observe(data.Duration.Seconds())

		if ok {
			metric.StoreReadMessageCount.Add(float64(data.Count))
		} else {
			metric.StoreReadRequestFailed.Inc()
		}

	})

	Driver.AddListener(MessageStored, func(i ...interface{}) {

		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*StoreEventData)
		if !ok {
			return
		}

		if !data.Success {
			metric.StoreStoreRequestFailed.Inc()
		}
		metric.StoreStoreDuration.Observe(data.Duration.Seconds())

	})

	Driver.AddListener(MessageDelete, func(i ...interface{}) {
		if len(i) != 1 {
			logger.Logger.Error("readStoreWriteToWriter error", zap.Any("i", i))
			return
		}
		data, ok := i[0].(*StoreEventData)
		if !ok {
			return
		}
		if !data.Success {
			metric.StoreDeleteRequestFailed.Inc()
		}
		metric.StoreDeleteDuration.Observe(data.Duration.Seconds())
	})

}

func TopicMessageStoredEventName(topic string) events.EventName {
	return WithEventPrefix(MessageStored, topic)
}
