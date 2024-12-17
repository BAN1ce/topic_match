package sub_topic

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/inner/broker/message_source"
	"github.com/BAN1ce/skyTree/inner/broker/store"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"time"
)

type Option func(topics *Topics)

// FullTopicName if it is a share topic,
// the share topic will like `$share/xxx/xxx`
type FullTopicName = string

type Topics struct {
	ctx    context.Context
	cancel context.CancelCauseFunc
	topic  map[FullTopicName]topic.Topic
}

func NewTopics(ctx context.Context, ops ...Option) *Topics {
	t := &Topics{
		topic: make(map[string]topic.Topic),
	}
	t.ctx, t.cancel = context.WithCancelCause(ctx)
	for _, op := range ops {
		op(t)
	}
	return t
}

// HandlePublishAck from a subtopic client
func (t *Topics) HandlePublishAck(topicName string, puback *packets.Puback) error {
	if puback.ReasonCode != packet.PublishAckSuccess {
		logger.Logger.Error().Str("topic", topicName).Str("client uid", pkg.GetClientUID(t.ctx)).Msg("handle publish Ack failed, reason code not success")
		return nil
	}

	if topicValue, ok := t.topic[topicName]; ok {
		if t, ok := topicValue.(topic.QoS1Subscriber); ok {
			return t.HandlePublishAck(puback)
		}

		logger.Logger.Error().Str("topic", topicName).Msg("handle publish Ack failed, handle type error not QoS1Subscriber")
		return errs.ErrTopicQoSNotSupport
	}

	logger.Logger.Warn().Str("topic", topicName).Msg("handle publish Ack failed, maybe topic not exists or handle type error")
	return errors.Join(errs.ErrTopicNotExistsInSubTopics, fmt.Errorf("topicName: %s", topicName))
}

func (t *Topics) HandlePublishRec(topicName string, pubrec *packets.Pubrec) error {
	if topicValue, ok := t.topic[topicName]; ok {
		if t, ok := topicValue.(topic.QoS2Subscriber); ok {
			return t.HandlePublishRec(pubrec)
		} else {
			return errors.Join(errs.ErrTopicQoSNotSupport, fmt.Errorf("topicName: %s", topicName))
		}

	}

	logger.Logger.Warn().Str("topic", topicName).Msg("handle publish Rec failed, topic not exists")

	return errors.Join(errs.ErrTopicNotExistsInSubTopics, fmt.Errorf("topicName: %s", topicName))
}

func (t *Topics) HandelPublishComp(topicName string, pubcomp *packets.Pubcomp) error {
	if topicValue, ok := t.topic[topicName]; ok {
		if t, ok := topicValue.(topic.QoS2Subscriber); ok {
			return t.HandlePublishComp(pubcomp)
		}
		logger.Logger.Warn().Str("topic", topicName).Msg("handle publish Comp failed, handle type error not QoS2Subscriber")
		return errs.ErrTopicQoSNotSupport

	}
	logger.Logger.Warn().Str("topic", topicName).Msg("handle publish Comp failed, topic not exists ")
	return errors.Join(errs.ErrTopicNotExistsInSubTopics, fmt.Errorf("topicName: %s", topicName))
}

func (t *Topics) AddTopic(topicName string, topic topic.Topic) (exists bool, err error) {
	if existTopic, ok := t.topic[topicName]; ok {
		exists = true
		if err = existTopic.Close(); err != nil {
			logger.Logger.Error().Err(err).Str("topic", topicName).Str("client uid", pkg.GetClientUID(t.ctx)).Msg("topic manager close close failed")
			return
		}
	}

	t.topic[topicName] = topic
	go func() {
		if err := topic.Start(t.ctx); err != nil {
			logger.Logger.Error().Err(err).Str("topic", topicName).Str("client uid", pkg.GetClientUID(t.ctx)).Msg("start topic failed")
		}
	}()
	return
}

func (t *Topics) DeleteTopic(topicName string) {
	if tt, ok := t.topic[topicName]; ok {
		if err := tt.Close(); err != nil {
			logger.Logger.Error().Err(err).Str("topic", topicName).Str("client uid", pkg.GetClientUID(t.ctx)).Msg("topic manager close topic failed")
		}
		delete(t.topic, topicName)
		return
	}
	logger.Logger.Warn().Str("topic", topicName).Msg("topics delete topic, topic not exists")

}

func (t *Topics) Close() error {
	// for no sub client
	if t == nil {
		return nil
	}
	t.cancel(fmt.Errorf("topics close"))
	return nil
}

func (t *Topics) GetUnfinishedMessage() map[string][]*packet.Message {
	var unfinishedMessage = make(map[string][]*packet.Message)
	for topicName, topicInstance := range t.topic {
		unfinishedMessage[topicName] = topicInstance.GetUnFinishedMessage()
	}
	return unfinishedMessage
}

func (t *Topics) Publish(topic string, message *packet.Message) error {
	if _, ok := t.topic[topic]; !ok {
		return errors.Join(errs.ErrTopicNotExistsInSubTopics, fmt.Errorf("topicName: %s", topic))
	}
	return t.topic[topic].Publish(message, nil)
}

func (t *Topics) Meta() []*topic.Meta {
	var meta = make([]*topic.Meta, 0)
	for _, to := range t.topic {
		meta = append(meta, to.Meta())
	}
	return meta
}

func CreateTopic(client client.PacketWriter, meta *topic.Meta) (topic.Topic, error) {
	var (
		topicName    = meta.Topic
		qos          = byte(meta.QoS)
		maxListenQoS = 0
	)
	// if subscribe qos is 0, the maxListenQoS is 2. The subscriber want to listen to the all level qos messages
	// if subscribe qos is 1 or 2, the maxListenQoS is 0, because the subscriber only wants to listen to the qos 0 messages
	// other qos level messages could read from store
	if qos == broker.QoS0 {
		maxListenQoS = 2
	}

	var stream = message_source.NewStreamSource(topicName, event.MessageEvent, maxListenQoS)

	logger.Logger.Debug().Str("client id", client.GetID()).Str("topic", topicName).Uint8("qos", qos).Msg("topic manager create topic")

	if !meta.Share {
		switch qos {
		case broker.QoS0:
			return NewQoS0(meta, client, stream), nil
		case broker.QoS1:
			return NewQoS1(meta, client, message_source.NewStoreSource(topicName, store.DefaultMessageStore, store.DefaultHandlePublishDoneEvent), stream, nil), nil
		case broker.QoS2:
			return NewQoS2(meta, client, message_source.NewStoreSource(topicName, store.DefaultMessageStore, store.DefaultHandlePublishDoneEvent), stream, nil), nil
		default:
			logger.Logger.Error().Str("client id", client.GetID()).Str("topic", topicName).Uint8("qos", qos).Msg("topic manager create topic failed, qos not support")
			return nil, fmt.Errorf("topic manager create topic failed, qos not support, topic: %s, qos: %d", topicName, qos)
		}
	}

	switch qos {
	case broker.QoS0:
		return NewQoS0(meta, client, store.DefaultShareMessageStore.GetShareGroupMessageSource(*meta.ShareTopic)), nil
	case broker.QoS1:
		return NewQoS1(meta, client, store.DefaultShareMessageStore.GetShareGroupMessageSource(*meta.ShareTopic), stream, nil), nil
	case broker.QoS2:
		return NewQoS2(meta, client, store.DefaultShareMessageStore.GetShareGroupMessageSource(*meta.ShareTopic), stream, nil), nil
	default:
		logger.Logger.Error().Str("client id", client.GetID()).Str("topic", topicName).Uint8("qos", qos).Msg("topic manager create topic failed, qos not support")
		return nil, fmt.Errorf("topic manager create topic failed, qos not support, topic: %s, qos: %d", topicName, qos)
	}

}

func CreateTopicFromSession(client client.PacketWriter, meta *topic.Meta, unfinished []*packet.Message) topic.Topic {
	var (
		topicName       = meta.Topic
		qos             = byte(meta.QoS)
		latestMessageID = meta.LatestMessageID
		ctx, cancel     = context.WithTimeout(context.TODO(), 10*time.Second)
		maxQoS          int
	)

	if qos != broker.QoS0 {
		maxQoS = 0
	}
	var stream = message_source.NewStreamSource(topicName, event.MessageEvent, maxQoS)

	defer cancel()

	switch byte(meta.QoS) {
	case broker.QoS0:
		return NewQoS0(meta, client, stream)

	case broker.QoS1:
		unfinished = FillUnfinishedMessage(ctx, unfinished, message_source.NewStoreSource(topicName, store.DefaultMessageStore, store.DefaultHandlePublishDoneEvent))
		return NewQoS1(meta, client, message_source.NewStoreSource(topicName, store.DefaultMessageStore, store.DefaultHandlePublishDoneEvent), stream, unfinished, QoS1WithLatestMessageID(latestMessageID))

	case broker.QoS2:
		unfinished = FillUnfinishedMessage(ctx, unfinished, message_source.NewStoreSource(topicName, store.DefaultMessageStore, store.DefaultHandlePublishDoneEvent))
		return NewQoS2(meta, client, message_source.NewStoreSource(topicName, store.DefaultMessageStore, store.DefaultHandlePublishDoneEvent), stream, unfinished, QoS2WithLatestMessageID(latestMessageID))

	default:
		logger.Logger.Error().Str("client id", client.GetID()).Str("topic", topicName).Uint8("qos", qos).Msg("topic manager create topic failed, qos not support")
		return nil
	}
}
