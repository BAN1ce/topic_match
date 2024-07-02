package event

import (
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"github.com/kataras/go-events"
)

var (
	Driver = events.New()
)

var GlobalEvent *Event

func Boot() {
	GlobalEvent = &Event{}
	createClientMetricListen()
}

type Event struct {
}

func (e *Event) CreatePublishEvent(topic string, handler func(...interface{})) {
	Driver.AddListener(ReceivedTopicPublishEventName(topic), handler)
}

func (e *Event) DeletePublishEvent(topic string, handler func(i ...interface{})) {
	Driver.RemoveListener(ReceivedTopicPublishEventName(topic), handler)
}

func (e *Event) EmitReceivedPublishDone(topic string, publish *packet.Message) {
	Driver.Emit(ReceivedTopicPublishEventName(topic), topic, publish)
}

func (e *Event) EmitMQTTPacket(eventName MQTTEventName, packet packets.Packet) {
	Driver.Emit(events.EventName(eventName), packet)
}
