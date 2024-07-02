package message_source

import (
	"context"
	broker2 "github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestNewStreamSource(t *testing.T) {
	var (
		mockCtrl            = gomock.NewController(t)
		mockPublishListener = broker.NewMockPublishListener(mockCtrl)
		topic               = "test"
		want                = &StreamSource{
			topic:           topic,
			publishListener: mockPublishListener,
		}
	)

	if got := NewStreamSource(topic, mockPublishListener); !reflect.DeepEqual(got, want) {
		t.Errorf("NewStreamSource() = %v, want %v", got, want)
	}

}

func TestStreamSource_Close(t *testing.T) {
	var (
		mockCtrl            = gomock.NewController(t)
		mockPublishListener = broker.NewMockPublishListener(mockCtrl)
		topic               = "test"
		ctx, cancel         = context.WithCancel(context.Background())
		messageChan         = make(chan *packet.Message, 10)
	)
	defer mockCtrl.Finish()
	mockPublishListener.EXPECT().DeletePublishEvent(topic, gomock.Any()).Return().AnyTimes()

	type fields struct {
		ctx             context.Context
		cancel          context.CancelFunc
		topic           string
		publishListener broker2.PublishListener
		messageChan     chan *packet.Message
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{

		{
			name:    "not listen",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "listen",
			fields: fields{
				ctx:             ctx,
				cancel:          cancel,
				messageChan:     messageChan,
				publishListener: mockPublishListener,
				topic:           topic,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StreamSource{
				ctx:             tt.fields.ctx,
				cancel:          tt.fields.cancel,
				topic:           tt.fields.topic,
				publishListener: tt.fields.publishListener,
				messageChan:     tt.fields.messageChan,
			}
			if err := s.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestListenMessage(t *testing.T) {
	var (
		mockCtrl            = gomock.NewController(t)
		mockPublishListener = broker.NewMockPublishListener(mockCtrl)
		topic               = "test"
		streamSource        = NewStreamSource(topic, mockPublishListener)
		ctx, cancel         = context.WithCancel(context.Background())
		messageChan         = make(chan *packet.Message, 10)
	)
	defer cancel()
	defer mockCtrl.Finish()

	mockPublishListener.EXPECT().CreatePublishEvent(topic, gomock.Any()).Return().Times(1)

	if err := streamSource.ListenMessage(ctx, messageChan); err != nil {
		t.Errorf("ListenMessage() error = %v", err)
	}

}

func TestHandler(t *testing.T) {
	var (
		mockCtrl            = gomock.NewController(t)
		mockPublishListener = broker.NewMockPublishListener(mockCtrl)
		topic               = "test"
		streamSource        = NewStreamSource(topic, mockPublishListener)
		ctx, cancel         = context.WithCancel(context.Background())
		messageChan         = make(chan *packet.Message, 10)
	)
	defer cancel()
	defer mockCtrl.Finish()

	mockPublishListener.EXPECT().CreatePublishEvent(topic, gomock.Any()).Return().Times(1)
	if err := streamSource.ListenMessage(ctx, messageChan); err != nil {
		t.Errorf("ListenMessage() error = %v", err)
	}

	streamSource.handler([]interface{}{
		"test",
		&packet.Message{
			PublishPacket: &packets.Publish{},
		},
	}...)

	if len(messageChan) != 1 {
		t.Errorf("handler() error = %v got=%d", "messageChan length is not 1", len(messageChan))
	}

}
