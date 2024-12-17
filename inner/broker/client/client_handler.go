package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/broker/client/sub_topic"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	client2 "github.com/BAN1ce/skyTree/pkg/broker/client"
	retain2 "github.com/BAN1ce/skyTree/pkg/broker/retain"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	topic2 "github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/rate"
	"github.com/BAN1ce/skyTree/pkg/retain"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"github.com/eclipse/paho.golang/packets"
)

var (
	pong = packets.NewControlPacket(packets.PINGRESP).Content.(*packets.Pingresp)
	// use the same ping resp packet
	pingResp = packets.NewControlPacket(packets.PINGRESP).Content.(*packets.Pingresp)
)

type InnerHandler struct {
	client *Client
}

func newClientHandler(client *Client) *InnerHandler {
	return &InnerHandler{
		client: client,
	}
}

func (i *InnerHandler) HandlePacket(ctx context.Context, packet *packets.ControlPacket, client *Client) error {
	var (
		err error
	)
	switch packet.FixedHeader.Type {
	case packets.CONNECT:
		connectPacket := packet.Content.(*packets.Connect)
		if err = i.HandleConnect(connectPacket); err != nil {
			logger.Logger.Warn().Err(err).Str("client", client.MetaString()).Msg("handle connect error")
		}

	case packets.PUBLISH:
		publishPacket := packet.Content.(*packets.Publish)
		return i.HandlePublish(publishPacket)

	case packets.SUBSCRIBE:
		subscribePacket := packet.Content.(*packets.Subscribe)
		if err = i.HandleSub(subscribePacket); err != nil {
			logger.Logger.Warn().Err(err).Str("client", client.MetaString()).Msg("handle subscribe error")
		}

	case packets.UNSUBSCRIBE:
		unsubscribePacket := packet.Content.(*packets.Unsubscribe)
		if err = i.HandleUnsub(unsubscribePacket); err != nil {
			logger.Logger.Warn().Err(err).Str("client", client.MetaString()).Msg("handle unsubscribe error")
		}

	case packets.PUBACK:
		pubAckPacket := packet.Content.(*packets.Puback)
		if err = i.HandlePubAck(pubAckPacket); err != nil {
			logger.Logger.Warn().Err(err).Str("client", client.MetaString()).Msg("handle pubAck error")
		}
	case packets.PUBREC:
		pubRecPacket := packet.Content.(*packets.Pubrec)
		err = i.HandlePubRec(pubRecPacket)

	case packets.PUBREL:
		err = i.HandlePubRel(packet.Content.(*packets.Pubrel))

	case packets.PUBCOMP:
		pubCompPacket := packet.Content.(*packets.Pubcomp)
		err = i.HandlePubComp(pubCompPacket)

	case packets.PINGREQ:
		return i.client.Write(&client2.WritePacket{
			Packet: pong,
		})

	case packets.DISCONNECT:
		err = i.client.Close()
	default:
		err = fmt.Errorf("unknown packet type = %d", packet.FixedHeader.Type)
		logger.Logger.Error().Err(err).Str("client", client.MetaString()).Msg("handle packet error")
	}
	return err
}

func (i *InnerHandler) HandleConnect(connectPacket *packets.Connect) error {
	i.client.mux.Lock()
	defer i.client.mux.Unlock()
	var (
		sessionConnectProp = session.NewConnectProperties(connectPacket.Properties)
		topicAliasMax      = uint16(1000)
	)
	i.client.connectProperties = sessionConnectProp
	if err := i.client.component.session.SetConnectProperties(sessionConnectProp); err != nil {
		return err
	}

	// if the clean start is false, recover the topic from the session
	if !connectPacket.CleanStart {
		i.recoverTopicFromSession()
	}

	var (
		windowSize = 0
	)

	if connectPacket.Properties.ReceiveMaximum != nil && *connectPacket.Properties.ReceiveMaximum > 0 {
		windowSize = int(*connectPacket.Properties.ReceiveMaximum)
	}

	i.client.publishBucket = rate.NewBucket(windowSize)

	if connectPacket.WillFlag {
		if err := i.client.setWill(&session.WillMessage{
			Topic:    connectPacket.WillTopic,
			QoS:      int(connectPacket.WillQOS),
			Property: &session.WillProperties{Properties: connectPacket.WillProperties},
			Retain:   false,
			Payload:  connectPacket.WillMessage,
			ClientID: i.client.GetID(),
		}); err != nil {
			return err
		}
	}

	var conAck = packets.NewControlPacket(packets.CONNACK).Content.(*packets.Connack)
	conAck.Properties = &packets.Properties{
		TopicAlias: &topicAliasMax,
	}
	conAck.ReasonCode = packets.ConnackSuccess
	return i.client.writePacket(&client2.WritePacket{
		Packet: conAck,
	})
}

// HandleSub handles the sub packet from the client
func (i *InnerHandler) HandleSub(subscribe *packets.Subscribe) error {
	c := i.client
	c.mux.Lock()

	var (
		subAck              = packets.NewControlPacket(packets.SUBACK).Content.(*packets.Suback)
		brokerTopics        = map[string]topic2.Topic{}
		topicRetainHandling = map[string]byte{}
		topicExists         = make(map[string]bool)
	)

	subAck.PacketID = subscribe.PacketID

	// create topic instance
	_, noShareSubscribePacket := topic2.SplitShareAndNoShare(subscribe)

	// handle simple sub
	for _, subOptions := range subscribe.Subscriptions {
		meta := topic2.NewMetaFromSubPacket(&subOptions, noShareSubscribePacket.Properties)
		t, err := sub_topic.CreateTopic(i.client, meta)
		if err != nil {
			subAck.Reasons = append(subAck.Reasons, 0x80)
			continue
		}
		brokerTopics[subOptions.Topic] = t

		i.client.getSession().CreateSubTopic(meta)

		if exists, err := c.topicManager.AddTopic(subOptions.Topic, t); err != nil {
			subAck.Reasons = append(subAck.Reasons, packet.ErrSubackCreateSubscriptionFailed)
			continue
		} else {
			topicExists[subOptions.Topic] = exists
		}

		if err := i.client.component.subCenter.CreateSub(i.client.ID, []packets.SubOptions{subOptions}); err != nil {
			subAck.Reasons = append(subAck.Reasons, packet.ErrSubackCreateSubscriptionFailed)
			continue
		}

		subAck.Reasons = append(subAck.Reasons, subOptions.QoS)

		topicRetainHandling[subOptions.Topic] = subOptions.RetainHandling
	}

	if err := i.client.writePacket(&client2.WritePacket{
		Packet: subAck,
	}); err != nil {
		return err
	}

	var (
		retainMessage = make(map[topic2.Topic]*packet.Message)
	)

	// client handle sub and create qos0,qos1,qos2 subOptions
	for topicName, t := range brokerTopics {

		if topicRetainHandling[topicName] == retain.HandingNoSendRetain {
			logger.Logger.Info().Str("topic", topicName).Msg("no send retain message")
			continue
		}

		if topicExists[topicName] && topicRetainHandling[topicName] == retain.HandingSendRetainOnlyNewConnection {
			logger.Logger.Info().Str("topic", topicName).Msg("doesn't send the retain message, because sending retain message only new subscription")
			continue
		}

		if utils.HasWildcard(topicName) {
			topics := c.component.subCenter.GetMatchTopicForWildcardTopic(topicName)
			for _, matchTopic := range topics {
				if message, ok := c.component.retain.GetRetainMessage(matchTopic); ok {
					logger.Logger.Debug().Str("topic", topicName).Msg("retain message publish")
					pub := packets.NewControlPacket(packets.PUBLISH).Content.(*packets.Publish)
					pub.Payload = message.Payload
					pub.Topic = topicName
					pub.Retain = true

					retainMessage[t] = &packet.Message{
						PublishPacket:  pub,
						Retain:         true,
						SubscribeTopic: topicName,
					}

				}
			}
			continue
		}

		if message, ok := c.component.retain.GetRetainMessage(topicName); ok {
			logger.Logger.Debug().Str("topic", topicName).Msg("retain message publish")
			pub := packets.NewControlPacket(packets.PUBLISH).Content.(*packets.Publish)
			pub.Payload = message.Payload
			pub.Topic = topicName
			pub.Retain = true

			retainMessage[t] = &packet.Message{
				PublishPacket: pub,
				Retain:        true,
			}

		}
	}
	c.mux.Unlock()

	// send retain messages
	for t, message := range retainMessage {
		if err := t.Publish(message, nil); err != nil {
			return errors.Join(err, fmt.Errorf("retain message publish"))
		}

	}
	return nil

}

func (i *InnerHandler) HandleUnsub(unsubscribe *packets.Unsubscribe) error {
	c := i.client
	c.mux.Lock()
	defer c.mux.Unlock()

	var (
		unsubAck = packets.NewControlPacket(packets.UNSUBACK).Content.(*packets.Unsuback)
	)
	unsubAck.PacketID = unsubscribe.PacketID

	for _, topicName := range unsubscribe.Topics {
		//  delete subscription first
		if err := c.component.subCenter.DeleteSub(i.client.ID, []string{topicName}); err != nil {
			logger.Logger.Error().Err(err).Str("client", c.ID).Str("topic", topicName).Msg("delete sub error")
			unsubAck.Reasons = append(unsubAck.Reasons, packets.UnsubackUnspecifiedError)
			continue
		}

		c.component.session.DeleteSubTopic(topicName)
		c.topicManager.DeleteTopic(topicName)
		unsubAck.Reasons = append(unsubAck.Reasons, packets.UnsubackSuccess)
	}

	return i.client.writePacket(&client2.WritePacket{
		Packet: unsubAck,
	})
}

func (i *InnerHandler) HandlePublish(publish *packets.Publish) error {
	var (
		topic     = publish.Topic
		response  packets.Packet
		messageID string
		err       error
		subTopics map[string]int32
		pubAck    = packet.NewPublishAck()
		pubrec    = packet.NewPublishRec()
	)

	pubAck.PacketID = publish.PacketID
	pubrec.PacketID = publish.PacketID

	response = pubAck
	if publish.QoS == broker.QoS2 {
		response = pubrec
	}

	// if the publication packet is duplicate, then response ack and return
	if publish.Duplicate {
		if err := i.client.Write(&client2.WritePacket{
			Packet: response,
		}); err != nil {
			logger.Logger.Warn().Err(err).Msg("write puback error")
		}
		logger.Logger.Info().Str("client", i.client.MetaString()).Str("topic", publish.Topic).Uint16("packetID", publish.PacketID).Msg("duplicate publish")
		return nil
	}

	defer func() {
		if response != nil {
			if err := i.client.Write(&client2.WritePacket{
				Packet: response,
			}); err != nil {
				logger.Logger.Warn().Err(err).Msg("write puback error")
			}
		}

		if publish.QoS == broker.QoS2 {
			// QoS2 need to wait for the pubrel packet
			return
		}

		if err != nil {
			// if err is not nil, then return, don't emit finish topic publish message event
			return
		}

		// emit finish topic publish message event for qos0, qos1
		for subTopic, qos := range subTopics {
			if len(subTopic) == 0 {
				continue
			}

			logger.Logger.Debug().Str("topic", subTopic).Int32("qos", qos).Msg("emit received publish event")
			// reset retain a flag, because the retain flag should not be true when doesn't it first subscribed
			publish.Retain = false
			event.MessageEvent.EmitFinishTopicPublishMessage(i.client.component.event, subTopic, &packet.Message{
				PublishPacket: publish,
				SendClientID:  i.client.GetID(),
				MessageID:     messageID,
			})
		}
	}()

	if err := i.handleTopicAlias(publish); err != nil {
		return err
	}

	if publish.QoS != broker.QoS2 {
		if publish.Retain {
			if len(publish.Payload) == 0 {
				if err := i.client.component.retain.DeleteRetainMessage(topic); err != nil {
					pubAck.ReasonCode = packets.PubackUnspecifiedError
				}
			} else if err := i.client.component.retain.PutRetainMessage(&retain2.Message{
				Topic:   topic,
				Payload: publish.Payload,
			}); err != nil {
				pubAck.ReasonCode = packets.PubackUnspecifiedError
			}
		}
	}

	if publish.QoS != broker.QoS2 {
		subTopics = i.client.component.subCenter.GetAllMatchTopics(publish.Topic)
	}

	// handle qos0, qos1, qos2

	switch publish.QoS {
	case broker.QoS0:

	case broker.QoS1:
		publish.Retain = false

		publishMessage := &packet.Message{
			SendClientID:  i.client.GetID(),
			PublishPacket: publish,
		}

		pubAck.ReasonCode = packets.PubackSuccess

		if len(subTopics) == 0 {
			logger.Logger.Info().Str("topic", publish.Topic).Msg("no sub topic")
			return nil
		}

		// messageStore message
		if messageID, err = i.client.component.messageStoreWrapper.StorePublishPacket(subTopics, publishMessage); err != nil {
			logger.Logger.Error().Err(err).Str("topic", topic).Str("messageID", messageID).Msg("messageStore publish packet for QoS1 error")
			pubAck.ReasonCode = packets.PubackUnspecifiedError
		}

		response = pubAck

	case broker.QoS2:

		pubrec.ReasonCode = packets.PubrecSuccess

		if publish.Duplicate {
			return nil
		}
		if !i.client.QoS2.StoreNotExists(publish) {
			pubrec.ReasonCode = packets.PubrecUnspecifiedError
		}
		response = pubrec
	}

	return nil
}

func (i *InnerHandler) HandlePubAck(pubAck *packets.Puback) error {
	c := i.client
	topicName := c.packetIdentifierIDTopic.GetTopic(pubAck.PacketID)
	if len(topicName) == 0 {
		logger.Logger.Warn().Str("client", c.MetaString()).Uint16("packetID", pubAck.PacketID).Msg("pubAck packetID not found store")
		return nil
	}
	if err := c.topicManager.HandlePublishAck(topicName, pubAck); err != nil {
		logger.Logger.Error().Err(err).Str("client", c.MetaString()).Any("meta", i.client.topicManager.Meta()).Msg("handle pubAck error")
		return err
	} else {
		c.publishBucket.PutToken()
	}
	return nil
}

func (i *InnerHandler) HandlePubRec(pubRec *packets.Pubrec) error {
	c := i.client
	topicName := c.packetIdentifierIDTopic.GetTopic(pubRec.PacketID)
	if len(topicName) == 0 {
		logger.Logger.Warn().Str("client", c.MetaString()).Uint16("packetID", pubRec.PacketID).Msg("pubRec packetID not found store")
		return nil
	}
	if err := c.topicManager.HandlePublishRec(topicName, pubRec); err != nil {
		logger.Logger.Error().Err(err).Str("client", c.MetaString()).Msg("handle pubRec error")
		return err
	}
	return nil
}

// HandlePubComp handles the pubComp packet from receiver
func (i *InnerHandler) HandlePubComp(pubRel *packets.Pubcomp) error {
	c := i.client
	topicName := c.packetIdentifierIDTopic.GetTopic(pubRel.PacketID)
	if len(topicName) == 0 {
		logger.Logger.Warn().Str("client", c.MetaString()).Uint16("packetID", pubRel.PacketID).Msg("pubComp packetID not found store")
		return nil
	}
	if err := c.topicManager.HandelPublishComp(topicName, pubRel); err != nil {
		logger.Logger.Error().Err(err).Str("client", c.MetaString()).Msg("handle pubComp error")
	} else {
		c.publishBucket.PutToken()
	}

	return nil
}

func (i *InnerHandler) HandlePubRel(publishRel *packets.Pubrel) (err error) {

	var (
		publishComp = packet.NewPublishComp()
		client      = i.client
		messageID   string
	)
	publishComp.PacketID = publishRel.PacketID
	publishPacket, ok := i.client.QoS2.Delete(publishRel)

	defer func() {
		err = client.Write(&client2.WritePacket{
			Packet: publishComp,
		})
		return
	}()

	if !ok {
		logger.Logger.Warn().Str("client", client.GetID()).Uint16("packetID", publishRel.PacketID).Msg("qos2 handle pubrel error, packet id not found, maybe deleted, because handled")
	}

	subTopics := i.client.component.subCenter.GetAllMatchTopics(publishPacket.Topic)
	if len(subTopics) == 0 {
		publishComp.ReasonCode = config.GetConfig().Broker.NoSubTopicResponse
		return
	}

	messageID, err = i.client.component.messageStoreWrapper.StorePublishPacket(subTopics, &packet.Message{
		SendClientID:  client.GetID(),
		PublishPacket: publishPacket,
	})

	if err != nil {
		logger.Logger.Error().Err(err).Msg("messageStore publish packet error")
		publishRel.ReasonCode = packets.PubackUnspecifiedError
		return
	}

	event.MessageEvent.EmitFinishTopicPublishMessage(i.client.component.event, publishPacket.Topic, &packet.Message{
		PublishPacket: publishPacket,
		MessageID:     messageID,
		SendClientID:  client.GetID(),
	})

	if !publishPacket.Retain {
		return
	}

	// if retain is true and payload is not empty, then create retain message id
	if len(publishPacket.Payload) > 0 {
		if err = i.client.component.retain.PutRetainMessage(&retain2.Message{
			Topic:   publishPacket.Topic,
			Payload: publishPacket.Payload,
		}); err != nil {
			logger.Logger.Error().Err(err).Str("topic", publishPacket.Topic).Msg("create retain message id error")
			publishRel.ReasonCode = packets.PubackUnspecifiedError
		}
	} else {
		// if retain is true and payload is empty, then delete retain message id
		if err = i.client.component.retain.DeleteRetainMessage(publishPacket.Topic); err != nil {
			logger.Logger.Error().Err(err).Str("topic", publishPacket.Topic).Msg("delete retain message id error")
			publishRel.ReasonCode = packets.PubackUnspecifiedError
		}
	}
	return

}

func (i *InnerHandler) HandlePing(_ *packets.Pingreq) error {
	i.client.RefreshAliveTime()
	err := i.client.component.monitorKeepAlive.SetClientAliveTime(i.client.GetID(), utils.NextAliveTime(int64(i.client.GetKeepAliveTime().Seconds())))
	if err != nil {
		logger.Logger.Error().Err(err).Str("client", i.client.GetID()).Msg("set client alive time error")
	}
	return i.client.Write(&client2.WritePacket{
		Packet: pingResp,
	})
}

func (i *InnerHandler) recoverTopicFromSession() {
	var (
		getSession = i.client.component.session
	)

	for _, topicMeta := range getSession.ReadSubTopics() {
		var (
			topicName = topicMeta.Topic
		)

		logger.Logger.Info().Str("id", i.client.GetID()).Str("topic", topicName).Msg("recover topic from getSession")
		unfinishedMessage := getSession.ReadTopicUnFinishedMessage(topicName)
		subTopic := sub_topic.CreateTopicFromSession(i.client, topicMeta, unfinishedMessage)
		if _, err := i.client.topicManager.AddTopic(topicName, subTopic); err != nil {
			logger.Logger.Error().Err(err).Str("topic", topicName).Msg("recover topic from getSession error")
		}
	}
}

func (i *InnerHandler) HandleDisconnect() error {
	//TODO implement me
	panic("implement me")
}

func (i *InnerHandler) HandleAuth() error {
	//TODO implement me
	panic("implement me")
}

func (p *InnerHandler) handleTopicAlias(packet *packets.Publish) error {
	client := p.client
	if packet.Properties == nil {
		return nil
	}
	if packet.Properties.TopicAlias == nil {
		return nil
	}
	if alias := *(packet.Properties.TopicAlias); alias != 0 {
		if packet.Topic != "" {
			client.SetClientTopicAlias(packet.Topic, *(packet.Properties.TopicAlias))
			return nil
		}
		// get the topic name from alias
		packet.Topic = client.GetClientTopicAlias(*(packet.Properties.TopicAlias))
		if packet.Topic == "" {
			return errs.ErrTopicAliasNotFound
		}
		// clear the topic alias, because the packet will use other where
		packet.Properties.TopicAlias = nil
		logger.Logger.Debug().Str("client", client.GetID()).Str("topic", packet.Topic).Msg("get topic from alias")
		return nil
	}
	return errs.ErrTopicAliasInvalid
}
