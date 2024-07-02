package message_source

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestReadStoreWriteToWriter(t *testing.T) {

	var ()

}

func TestNewStoreSource(t *testing.T) {
	var (
		mockCtrl              = gomock.NewController(t)
		mockMessageStore      = broker.NewMockMessageStore(mockCtrl)
		mockMessageStoreEvent = broker.NewMockMessageStoreEvent(mockCtrl)
		topic                 = "test"
		result                = &StoreSource{
			topic:       topic,
			messageChan: make(chan *packet.Message, 10),
			storeEvent:  mockMessageStoreEvent,
			store:       mockMessageStore,
		}
	)
	defer mockCtrl.Finish()

	if got := NewStoreSource(topic, mockMessageStore, mockMessageStoreEvent); reflect.DeepEqual(got, result) {
		t.Errorf("NewStoreSource() = %v, want %v", got, result)
	}
}

func TestNextMessages(t *testing.T) {
	var (
		mockCtrl              = gomock.NewController(t)
		mockMessageStore      = broker.NewMockMessageStore(mockCtrl)
		mockMessageStoreEvent = broker.NewMockMessageStoreEvent(mockCtrl)
		topic                 = "test"
		storeSource           = NewStoreSource(topic, mockMessageStore, mockMessageStoreEvent)
		messages              []*packet.Message
	)

	for i := 0; i < 10; i++ {
		messages = append(messages, &packet.Message{
			MessageID: fmt.Sprintf("%d", i),
		})
	}

	mockMessageStore.EXPECT().ReadTopicMessagesByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(messages, nil).AnyTimes()

	got, err := storeSource.NextMessages(context.TODO(), 10, "0", true)
	if err != nil {
		t.Errorf("NextMessages() error = %v", err)
	}

	if !reflect.DeepEqual(got, messages) {
		t.Errorf("NextMessages() = %v, want %v", got, messages)
	}

	for i := 0; i < 10; i++ {
		if got[i].MessageID != messages[i].MessageID {
			t.Errorf("NextMessages() = %v, want %v", got[i].MessageID, messages[i].MessageID)
		}
	}
}

func TestNextMessageWithWaitingEvent(t *testing.T) {
	var (
		mockCtrl              = gomock.NewController(t)
		mockMessageStore      = broker.NewMockMessageStore(mockCtrl)
		mockMessageStoreEvent = broker.NewMockMessageStoreEvent(mockCtrl)
		topic                 = "test"
		storeSource           = NewStoreSource(topic, mockMessageStore, mockMessageStoreEvent)
		messages              []*packet.Message
	)

	for i := 0; i < 10; i++ {
		messages = append(messages, &packet.Message{
			MessageID: fmt.Sprintf("%d", i),
		})
	}

	mockMessageStore.EXPECT().ReadTopicMessagesByID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(messages, nil).AnyTimes()

	mockMessageStoreEvent.EXPECT().CreateListenMessageStoreEvent(gomock.Any(), gomock.Any()).DoAndReturn(func(topic string, f func(i ...interface{})) {
		f("test", "10")

	}).AnyTimes()
	mockMessageStoreEvent.EXPECT().DeleteListenMessageStoreEvent(gomock.Any(), gomock.Any()).AnyTimes()

	got, err := storeSource.NextMessages(context.TODO(), 10, "", true)
	if err != nil {
		t.Errorf("NextMessages() error = %v", err)
	}

	if !reflect.DeepEqual(got, messages) {
		t.Errorf("NextMessages() = %v, want %v", got, messages)
	}

	for i := 0; i < 10; i++ {
		if got[i].MessageID != messages[i].MessageID {
			t.Errorf("NextMessages() = %v, want %v", got[i].MessageID, messages[i].MessageID)
		}
	}

	got, err = storeSource.NextMessages(context.TODO(), 10, "0", true)
	if err != nil {
		t.Errorf("NextMessages() error = %v", err)
	}

	if !reflect.DeepEqual(got, messages) {
		t.Errorf("NextMessages() = %v, want %v", got, messages)
	}

	for i := 0; i < 10; i++ {
		if got[i].MessageID != messages[i].MessageID {
			t.Errorf("NextMessages() = %v, want %v", got[i].MessageID, messages[i].MessageID)
		}
	}

}
