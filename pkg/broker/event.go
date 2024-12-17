package broker

import event2 "github.com/BAN1ce/skyTree/inner/event"

type HandlePublishDoneEvent interface {
	CreateListenPublishEvent(topic string, handler func(i ...interface{}))
	DeleteListenPublishEvent(topic string, handler func(i ...interface{}))
	CreateListenPublishEventOnce(topic string, handler func(...interface{}))
}

type MessageStoreEvent interface {
	CreateListenMessageStoreEvent(topic string, handler func(...interface{}))
	CreateOnceListenMessageStoreEvent(topic string, handler func(...interface{}))
	DeleteListenMessageStoreEvent(topic string, handler func(i ...interface{}))
	EmitStored(*event2.StoreEventData)
	EmitRead(e *event2.StoreEventData)
	EmitDelete(e *event2.StoreEventData)
}
