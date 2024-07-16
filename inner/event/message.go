package event

import (
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/utils"
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

// DeletePublishEvent  deletes the event listener for the given topic.
func (e *Message) DeletePublishEvent(topic string, handler func(i ...interface{})) {
	eventDriver.RemoveListener(ReceivedTopicPublishEventName(topic), handler)
}

// EmitReceivedPublishDone emits the event when a message is published to the topic.
// Different QoS level triggers by different MQTT packet.
// Like QoS 1 triggers by the PUBACK packet.
// QoS 2 triggers by the PUBREL packet.
// QoS 0 triggers by the PUBLISH packet
func (e *Message) EmitReceivedPublishDone(topic string, publish *packet.Message) {
	eventDriver.Emit(ReceivedTopicPublishEventName(topic), topic, publish)
}

const (
	ReceivedTopicPublishEvent = "event.received.topic.publish"
)

func ReceivedTopicPublishEventName(topic string) string {
	return utils.WithEventPrefix(ReceivedTopicPublishEvent, topic)
}
