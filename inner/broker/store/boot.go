package store

import (
	"context"
	"github.com/BAN1ce/skyTree/inner/broker/share"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/message/serializer"
)

var (
	DefaultMessageStore      broker.TopicMessageStore
	DefaultMessageStoreEvent broker.MessageStoreEvent
	DefaultSerializerVersion serializer.SerialVersion

	DefaultShareMessageStore *share.TopicMessageSourceFactory
)

func Boot(store broker.TopicMessageStore, event broker.MessageStoreEvent) {
	DefaultMessageStore = store
	DefaultMessageStoreEvent = event
	DefaultSerializerVersion = serializer.ProtoBufVersion

	DefaultShareMessageStore = share.NewTopicMessageSourceFactory(context.TODO(), store, event)
}
