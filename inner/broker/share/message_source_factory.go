package share

import (
	"context"
	"github.com/BAN1ce/skyTree/inner/broker/message_source"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	topic2 "github.com/BAN1ce/skyTree/pkg/broker/topic"
	"go.uber.org/zap"
	"sync"
)

type GroupName = string

type TopicName = string

type GroupMessageSource = map[GroupName]*MessageSource

// TopicMessageSourceFactory is the factory of the message source
// Every share group for one topic has a message source, in a share group client use the same message source,
// So we need a factory to create the message source.
type TopicMessageSourceFactory struct {
	ctx              context.Context
	mux              sync.RWMutex
	ShareTopicSource map[TopicName]GroupMessageSource
	store            broker.TopicMessageStore
	storeEvent       broker.MessageStoreEvent
}

func NewTopicMessageSourceFactory(ctx context.Context, store broker.TopicMessageStore, storeEvent broker.MessageStoreEvent) *TopicMessageSourceFactory {
	return &TopicMessageSourceFactory{
		ctx:              ctx,
		ShareTopicSource: make(map[TopicName]GroupMessageSource),
		store:            store,
		storeEvent:       storeEvent,
	}
}

// GetShareGroupMessageSource get the message source of the share group for the topic
// If the message source does not exist, create a new message source.
func (m *TopicMessageSourceFactory) GetShareGroupMessageSource(shareTopic topic2.ShareTopic) broker.MessageSource {
	m.mux.Lock()
	defer m.mux.Unlock()

	var (
		topic          = shareTopic.Topic
		group          = shareTopic.ShareGroup
		fullShareTopic = shareTopic.FullTopic
	)

	if _, ok := m.ShareTopicSource[topic]; !ok {
		m.ShareTopicSource[topic] = make(GroupMessageSource)
	}
	if _, ok := m.ShareTopicSource[topic][group]; !ok {
		topicMessageSource := message_source.NewStoreSource(shareTopic.Topic, m.store, m.storeEvent)
		m.ShareTopicSource[topic][group] = newMessageSource(fullShareTopic, topic, topicMessageSource)
	}
	return m.ShareTopicSource[topic][group]
}

func (m *TopicMessageSourceFactory) DeleteShareGroupMessageSource(shareTopic topic2.ShareTopic) {
	m.mux.Lock()
	defer m.mux.Unlock()

	var (
		topic = shareTopic.Topic
		group = shareTopic.ShareGroup
	)

	if _, ok := m.ShareTopicSource[topic]; !ok {
		// TODO add metric for this maybe better
		return
	}
	if _, ok := m.ShareTopicSource[topic][group]; !ok {
		// TODO add metric for this maybe better
		return
	}
	if err := m.ShareTopicSource[topic][group].Close(); err != nil {
		logger.Logger.Error("close share group message source error", zap.Error(err))
	}
	delete(m.ShareTopicSource[topic], group)

}
