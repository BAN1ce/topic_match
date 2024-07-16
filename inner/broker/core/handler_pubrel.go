package core

import (
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/logger"
	client2 "github.com/BAN1ce/skyTree/pkg/broker/client"
	packet2 "github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"go.uber.org/zap"
)

type PublishRel struct {
}

func NewPublishRel() *PublishRel {
	return &PublishRel{}
}

func (p *PublishRel) Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error) {
	var (
		publishRel  = rawPacket.Content.(*packets.Pubrel)
		publishComp = packet2.NewPublishComp()
		messageID   string
	)
	publishComp.PacketID = publishRel.PacketID
	publishPacket, ok := client.QoS2.Delete(publishRel)

	if !ok {
		logger.Logger.Warn("qos2 handle pubrel error, packet id not found, maybe deleted, because handled", zap.String("clientID", client.GetID()), zap.Uint16("packetID", publishRel.PacketID))
		return client.WritePacket(client2.NewWritePacket(publishComp))
	} else {
		// Emit event
		event.MessageEvent.EmitReceivedPublishDone(publishPacket.Topic, &packet2.Message{
			PublishPacket: publishPacket,
		})
	}

	subTopics := broker.subTree.MatchTopic(publishPacket.Topic)
	if len(subTopics) == 0 {
		publishComp.ReasonCode = config.GetConfig().Broker.NoSubTopicResponse
	} else {
		messageID, err = broker.messageStore.StorePublishPacket(subTopics, &packet2.Message{
			ClientID:      client.GetID(),
			PublishPacket: publishPacket,
		})
		if err != nil {
			logger.Logger.Error("messageStore publish packet error", zap.Error(err))
			publishRel.ReasonCode = packets.PubackUnspecifiedError
		}
		if !publishPacket.Retain {
			return client.WritePacket(client2.NewWritePacket(publishComp))
		}
	}
	// if retain is true and payload is not empty, then create retain message id
	if len(publishPacket.Payload) > 0 {
		for topic := range subTopics {
			if err = broker.state.CreateRetainMessageID(topic, messageID); err != nil {
				logger.Logger.Error("create retain message id error", zap.Error(err), zap.String("topic", topic))
				publishRel.ReasonCode = packets.PubackUnspecifiedError
			}
		}
	} else {
		// if retain is true and payload is empty, then delete retain message id
		if err = broker.state.DeleteRetainMessageID(publishPacket.Topic); err != nil {
			logger.Logger.Error("delete retain message id error", zap.Error(err))
			publishRel.ReasonCode = packets.PubackUnspecifiedError
		}
	}
	return client.WritePacket(&client2.WritePacket{
		Packet: publishComp,
	})
}
