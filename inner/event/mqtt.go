package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/eclipse/paho.golang/packets"
)

type MQTTEventName = string

const (
	MQTTReceivedPacketEvent MQTTEventName = "event.received.mqtt.packet"
	MQTTSendPacketEvent     MQTTEventName = "event.send.mqtt.packet"
)

var (
	MQTTEvent = newMQTT()
)

var (
	packetTypeString = []string{
		"",
		"CONNECT",
		"CONNACK",
		"PUBLISH",
		"PUBACK",
		"PUBREC",
		"PUBREL",
		"PUBCOMP",
		"SUBSCRIBE",
		"SUBACK",
		"UNSUBSCRIBE",
		"UNSUBACK",
		"PINGREQ",
		"PINGRESP",
		"DISCONNECT",
		"AUTH",
	}
)

type MQTT struct {
}

func newMQTT() *MQTT {
	m := &MQTT{}

	return m

}

func (e *MQTT) EmitReceivedMQTTPacketEvent(packetType byte) {
	eventDriver.Emit(MQTTReceivedPacketEvent, packetType)
}

func (e *MQTT) EmitSendMQTTPacketEvent(packetType byte) {
	logger.Logger.Debug().Str("packet type", packetTypeString[packetType]).Msg("write packet")
	eventDriver.Emit(MQTTSendPacketEvent, packetType)
}

func (e *MQTT) addDefaultListener() {
	e.registerReceivedMetric()
	e.registerSendMetric()
}

func (e *MQTT) registerReceivedMetric() {
	eventDriver.AddListener(MQTTReceivedPacketEvent, func(i ...interface{}) {
		if len(i) != 1 {
			logger.Logger.Error().Msg("invalid parameter")
			return
		}

		event, ok := i[0].(byte)
		if !ok {
			return
		}

		switch event {
		case packets.CONNECT:
			metric.ReceivedConnect.Inc()

		case packets.DISCONNECT:
			metric.ReceivedDisconnect.Inc()

		case packets.PINGREQ:
			metric.ReceivedPing.Inc()

		case packets.SUBSCRIBE:
			metric.ReceivedSubscription.Inc()

		case packets.UNSUBSCRIBE:
			metric.ReceivedUnsubscription.Inc()

		case packets.PUBLISH:
			metric.ReceivedPublish.Inc()

		case packets.PUBACK:
			metric.ReceivedPublishAck.Inc()

		case packets.PUBCOMP:
			metric.ReceivedPubComp.Inc()

		case packets.PUBREC:
			metric.ReceivedPubRec.Inc()

		case packets.PUBREL:
			metric.ReceivedPubRel.Inc()

		case packets.AUTH:
			metric.ReceivedAuth.Inc()
		default:
			logger.Logger.Error().Any("packet", event).Msg("unknown packet type")
		}

	})
}

func (e *MQTT) registerSendMetric() {
	eventDriver.AddListener(MQTTSendPacketEvent, func(i ...interface{}) {
		if len(i) != 1 {
			logger.Logger.Error().Msg("invalid parameter")
			return
		}

		event, ok := i[0].(byte)
		if !ok {
			return
		}

		switch event {

		case packets.DISCONNECT:
			metric.SendDisconnect.Inc()

		case packets.CONNACK:
			metric.SendConnectAck.Inc()

		case packets.PINGRESP:
			metric.SendPong.Inc()

		case packets.SUBACK:
			metric.SendSubscriptionAck.Inc()

		case packets.UNSUBACK:
			metric.SendUnsubscriptionAck.Inc()

		case packets.PUBLISH:
			metric.SendPublish.Inc()

		case packets.PUBACK:
			metric.SendPublishAck.Inc()

		case packets.PUBCOMP:
			metric.SendPubComp.Inc()

		case packets.PUBREC:
			metric.SendPubRec.Inc()

		case packets.PUBREL:
			metric.SendPubRel.Inc()

		case packets.AUTH:
			metric.SendAuthAck.Inc()
		default:
			logger.Logger.Error().Any("packet", event).Msg("unknown packet type")
		}
	})
}
