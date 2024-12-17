package event

import (
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/utils"
)

const (
	FinishTopicPublishMessageEvent = "event.finish.topic.publish.message"
)

var (
	MessageEvent = newMessageEvent()
)

type Message struct {
}

func newMessageEvent() *Message {
	return &Message{}
}

// CreateListenPublishEvent creates a new event listener for the given topic.
// The handler will be called when a new message is published to the topic.
// A message only trigger the handler once.
// NO matter which the QoS level is.
func (e *Message) CreateListenPublishEvent(topic string, handler func(...interface{})) {
	eventDriver.AddListener(ReceivedTopicPublishEventName(topic), handler)
}

func (e *Message) CreateListenPublishEventOnce(topic string, handler func(...interface{})) {
	eventDriver.AddListenerOnce(ReceivedTopicPublishEventName(topic), handler)
}

// DeleteListenPublishEvent DeletePublishEvent  deletes the event listener for the given topic.
func (e *Message) DeleteListenPublishEvent(topic string, handler func(i ...interface{})) {
	eventDriver.RemoveListener(ReceivedTopicPublishEventName(topic), handler)
}

// EmitFinishTopicPublishMessage emits the event when a message is published to the topic.
func (e *Message) EmitFinishTopicPublishMessage(eventDriver Driver, topic string, publish *packet.Message) {
	eventDriver.Emit(ReceivedTopicPublishEventName(topic), topic, publish)
}

func ReceivedTopicPublishEventName(topic string) string {
	return utils.WithEventPrefix(FinishTopicPublishMessageEvent, topic)
}
