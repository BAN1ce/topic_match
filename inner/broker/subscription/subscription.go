package subscription

import (
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/eclipse/paho.golang/packets"
	"time"
)

type Wrapper struct {
	subscribeCenter broker.SubCenter
}

func NewWrapper(center broker.SubCenter) *Wrapper {
	return &Wrapper{
		subscribeCenter: center,
	}
}

func (w *Wrapper) CreateSub(clientID string, topics []packets.SubOptions) error {
	var (
		start = time.Now()
		err   = w.subscribeCenter.CreateSub(clientID, topics)
	)
	event.SubtreeEvent.Emit(event.SubtreeWriteEvent, &event.SubtreeEventData{
		Success:  err == nil,
		Duration: time.Since(start),
	})
	return err
}

func (w *Wrapper) DeleteSub(clientID string, topics []string) error {
	var (
		start = time.Now()
		err   = w.subscribeCenter.DeleteSub(clientID, topics)
	)
	event.SubtreeEvent.Emit(event.SubtreeDeleteEvent, &event.SubtreeEventData{
		Success:  err == nil,
		Duration: time.Since(start),
	})
	return err
}

func (w *Wrapper) Match(topic string) (clientIDQos map[string]int32) {
	var (
		start  = time.Now()
		result = w.subscribeCenter.Match(topic)
	)
	event.SubtreeEvent.Emit(event.SubtreeReadEvent, &event.SubtreeEventData{
		Success:      true,
		Duration:     time.Since(start),
		NoSubscriber: len(result) == 0,
	})
	return result
}

func (w *Wrapper) DeleteClient(clientID string) error {
	var (
		start = time.Now()
		err   = w.subscribeCenter.DeleteClient(clientID)
	)
	event.SubtreeEvent.Emit(event.SubtreeDeleteEvent, &event.SubtreeEventData{
		Success:  err == nil,
		Duration: time.Since(start),
	})
	return err
}

func (w *Wrapper) MatchTopic(topic string) (topics map[string]int32) {
	var (
		start  = time.Now()
		result = w.subscribeCenter.MatchTopic(topic)
	)
	event.SubtreeEvent.Emit(event.SubtreeReadEvent, &event.SubtreeEventData{
		Success:      true,
		Duration:     time.Since(start),
		NoSubscriber: len(result) == 0,
	})
	return result
}
