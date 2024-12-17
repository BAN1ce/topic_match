package store

import (
	"github.com/BAN1ce/skyTree/inner/broker/share"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/message/serializer"
)

var (
	DefaultMessageStore broker.TopicMessageStore

	DefaultHandlePublishDoneEvent broker.HandlePublishDoneEvent

	DefaultSerializerVersion serializer.SerialVersion

	DefaultShareMessageStore *share.TopicMessageSourceFactory
)

func SetDefault(store broker.TopicMessageStore, event broker.HandlePublishDoneEvent, serializerVersion serializer.SerialVersion) {
	DefaultMessageStore = store
	DefaultHandlePublishDoneEvent = event
	DefaultSerializerVersion = serializer.ProtoBufVersion
	DefaultShareMessageStore = share.NewTopicMessageSourceFactory(store, event)
}
