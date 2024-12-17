package sub_center

import (
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/eclipse/paho.golang/packets"
	"time"
)

type SubCenterWrapper struct {
	subscribeCenter broker.SubCenter
}

func NewSubCenterWrapper(center broker.SubCenter) *SubCenterWrapper {
	return &SubCenterWrapper{
		subscribeCenter: center,
	}
}

func (w *SubCenterWrapper) CreateSub(clientID string, topics []packets.SubOptions) error {
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

func (w *SubCenterWrapper) DeleteSub(clientID string, topics []string) error {
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

func (w *SubCenterWrapper) Match(topic string) (clientIDQos map[string]int32) {
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

func (w *SubCenterWrapper) DeleteClient(clientID string) error {
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

func (w *SubCenterWrapper) GetAllMatchTopics(topic string) (topics map[string]int32) {
	var (
		start  = time.Now()
		result = w.subscribeCenter.GetAllMatchTopics(topic)
	)
	event.SubtreeEvent.Emit(event.SubtreeReadEvent, &event.SubtreeEventData{
		Success:      true,
		Duration:     time.Since(start),
		NoSubscriber: len(result) == 0,
	})
	return result
}

func (w *SubCenterWrapper) GetMatchTopicForWildcardTopic(wildTopic string) []string {
	var (
		result = w.subscribeCenter.GetMatchTopicForWildcardTopic(wildTopic)
	)
	return result
}
