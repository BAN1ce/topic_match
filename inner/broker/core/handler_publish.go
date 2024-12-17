package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/BAN1ce/skyTree/logger"
	broker2 "github.com/BAN1ce/skyTree/pkg/broker"
	client2 "github.com/BAN1ce/skyTree/pkg/broker/client"
	packet2 "github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
)

type PublishHandler struct {
}

func NewPublishHandler() *PublishHandler {
	return &PublishHandler{}
}

func (p *PublishHandler) Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) error {
	//var (
	//	packet     = rawPacket.Content.(*packets.Publish)
	//	err        error
	//	reasonCode byte
	//)
	//
	//// response puback or pubrec at the end of this function
	//defer func() {
	//	p.response(client, reasonCode, packet)
	//}()
	//
	//// if a packet is duplicate, then ignore it
	//if packet.Duplicate {
	//	return nil
	//}
	//
	//// handle topic alias, if topic alias is not 0, then use topic alias
	//// if topic alias is 0, then use topic
	//if err = p.handleTopicAlias(packet, client); err != nil {
	//	logger.Logger.Error().Err(err).Msg("handle topic alias error")
	//	reasonCode = packets.PubackUnspecifiedError
	//	p.response(client, reasonCode, packet)
	//	return err
	//}

	return nil
}

func (p *PublishHandler) response(client *client.Client, reasonCode byte, publish *packets.Publish) {
	switch publish.QoS {
	case broker2.QoS1:
		pubAck := packet2.NewPublishAck()
		pubAck.ReasonCode = reasonCode
		pubAck.PacketID = publish.PacketID
		if err := client.Write(client2.NewWritePacket(pubAck)); err != nil {
			logger.Logger.Warn().Err(err).Msg("write puback error")
		}

	case broker2.QoS2:
		pubRec := packet2.NewPublishRec()
		pubRec.ReasonCode = reasonCode
		pubRec.PacketID = publish.PacketID
		if err := client.Write(client2.NewWritePacket(pubRec)); err != nil {
			logger.Logger.Warn().Err(err).Msg("write pubrec error")
		}
	case broker2.QoS0:
		// do nothing
	default:
		logger.Logger.Error().Uint8("qos", publish.QoS).Msg("wrong qos")
	}

}
