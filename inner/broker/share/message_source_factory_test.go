package share

import (
	"context"
	topic2 "github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestNewTopicMessageSourceFactory(t *testing.T) {
	var (
		mockCtrl                  = gomock.NewController(t)
		mockTopicMessageSource    = broker.NewMockTopicMessageStore(mockCtrl)
		topicMessageSourceFactory = broker.NewMockMessageStoreEvent(mockCtrl)
		ctx, cancel               = context.WithCancel(context.Background())
		wantGot                   = &TopicMessageSourceFactory{
			ctx:              ctx,
			ShareTopicSource: make(map[TopicName]GroupMessageSource),
			store:            mockTopicMessageSource,
			storeEvent:       topicMessageSourceFactory,
		}
	)

	defer mockCtrl.Finish()
	defer cancel()

	result := NewTopicMessageSourceFactory(ctx, mockTopicMessageSource, topicMessageSourceFactory)

	if !reflect.DeepEqual(result, wantGot) {
		t.Errorf("NewTopicMessageSourceFactory() = %v, want %v", result, wantGot)
	}
}

func TestTopicMessageSourceFactory_GetShareGroupMessageSource_DeleteShareGroupMessageSource(t *testing.T) {
	var (
		mockCtrl                  = gomock.NewController(t)
		mockTopicMessageSource    = broker.NewMockTopicMessageStore(mockCtrl)
		topicMessageSourceFactory = broker.NewMockMessageStoreEvent(mockCtrl)
		ctx, cancel               = context.WithCancel(context.Background())
		factory                   = NewTopicMessageSourceFactory(ctx, mockTopicMessageSource, topicMessageSourceFactory)
		shareTopic                = topic2.ShareTopic{
			Topic:      "topic",
			ShareGroup: "group",
			FullTopic:  "$share/group/topic",
		}
	)
	defer mockCtrl.Finish()
	defer cancel()

	got1 := factory.GetShareGroupMessageSource(shareTopic)

	if got1 == nil {
		t.Errorf("GetShareGroupMessageSource() = %v, want %v", got1, shareTopic)
	}

	got2 := factory.GetShareGroupMessageSource(shareTopic)

	if got2 == nil {
		t.Errorf("GetShareGroupMessageSource() = %v, want %v", got2, shareTopic)
	}

	if got1 != got2 {
		t.Errorf("GetShareGroupMessageSource() = %v, want %v", got1, got2)
	}

	factory.DeleteShareGroupMessageSource(topic2.ShareTopic{
		Topic:      "topic",
		ShareGroup: "group1",
	})

	factory.DeleteShareGroupMessageSource(topic2.ShareTopic{
		Topic: "topic1",
	})

	factory.DeleteShareGroupMessageSource(shareTopic)

	got3 := factory.GetShareGroupMessageSource(shareTopic)

	if got3 == nil {
		t.Errorf("GetShareGroupMessageSource() = %v, want %v", got3, shareTopic)
	}

	// DeleteShareGroupMessageSource() should delete the share group message source,
	// so got3 should not be equal to got1
	if got3 == got1 {
		t.Errorf("GetShareGroupMessageSource() = %v, want %v", got3, got1)
	}

}
