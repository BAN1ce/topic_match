package sub_topic

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	broker2 "github.com/BAN1ce/skyTree/pkg/mocks/broker"
	client2 "github.com/BAN1ce/skyTree/pkg/mocks/client"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
	"time"
)

func TestQoS0_GetUnfinishedMessage(t *testing.T) {
	type fields struct {
		ctx           context.Context
		cancel        context.CancelFunc
		topic         string
		messageSource broker.MessageSource
		client        *Client
		meta          *topic.Meta
	}
	tests := []struct {
		name   string
		fields fields
		want   []*packet.Message
	}{
		{
			name:   "qos0 unfinished message",
			fields: fields{},
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS0{
				ctx:           tt.fields.ctx,
				cancel:        tt.fields.cancel,
				topic:         tt.fields.topic,
				messageSource: tt.fields.messageSource,
				client:        tt.fields.client,
				meta:          tt.fields.meta,
			}
			if got := q.GetUnFinishedMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUnfinishedMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewQoS0(t *testing.T) {
	type args struct {
		meta          *topic.Meta
		writer        client.PacketWriter
		messageSource broker.MessageSource
	}
	var (
		topicMeta = &topic.Meta{
			Topic: "test",
		}
		mockCtrl = gomock.NewController(t)
	)
	defer mockCtrl.Finish()

	mockWriter := client2.NewMockPacketWriter(mockCtrl)
	mockMessageSource := broker2.NewMockMessageSource(mockCtrl)
	tests := []struct {
		name string
		args args
		want *QoS0
	}{
		{
			name: "new qos0",
			args: args{
				meta:          topicMeta,
				writer:        mockWriter,
				messageSource: mockMessageSource,
			},
			want: &QoS0{
				topic:         topicMeta.Topic,
				meta:          topicMeta,
				messageSource: mockMessageSource,
				client:        NewClient(mockWriter, topicMeta),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQoS0(tt.args.meta, tt.args.writer, tt.args.messageSource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQoS0() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQoS0_GetID(t *testing.T) {
	type fields struct {
		ctx           context.Context
		cancel        context.CancelFunc
		topic         string
		messageSource broker.MessageSource
		client        *Client
		meta          *topic.Meta
	}
	var (
		mockCtrl = gomock.NewController(t)
		result   = "1"
	)
	defer mockCtrl.Finish()

	mockWriter := client2.NewMockPacketWriter(mockCtrl)

	mockWriter.EXPECT().GetID().Return(result).AnyTimes()

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "get id",
			fields: fields{
				client: NewClient(mockWriter, nil),
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS0{
				ctx:           tt.fields.ctx,
				cancel:        tt.fields.cancel,
				topic:         tt.fields.topic,
				messageSource: tt.fields.messageSource,
				client:        tt.fields.client,
				meta:          tt.fields.meta,
			}
			if got := q.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQoS0_Publish(t *testing.T) {
	type fields struct {
		ctx           context.Context
		cancel        context.CancelFunc
		topic         string
		messageSource broker.MessageSource
		client        *Client
		meta          *topic.Meta
		messageChan   chan *packet.Message
	}
	type args struct {
		publish *packet.Message
	}
	var (
		ctx, cancel   = context.WithTimeout(context.TODO(), 3*time.Second)
		mockCtrl      = gomock.NewController(t)
		mockWriter    = client2.NewMockPacketWriter(mockCtrl)
		publishPacket = packets.Publish{
			Payload: []byte("ssss"),
			Topic:   "sss",
		}
		metaTopic = &topic.Meta{}

		mockMessageSource = broker2.NewMockMessageSource(mockCtrl)
		listenChan        = make(chan *packet.Message, 10)
	)
	defer mockCtrl.Finish()

	mockWriter.EXPECT().NextPacketID().AnyTimes().Return(uint16(1))

	mockWriter.EXPECT().WritePacket(gomock.Any()).AnyTimes().DoAndReturn(func(publish *client.WritePacket) error {
		if reflect.DeepEqual(publish, &client.WritePacket{
			Packet: &publishPacket,
		}) {
			return nil
		}
		return fmt.Errorf("not match")
	})

	mockWriter.EXPECT().GetID().AnyTimes().Return("1")

	mockMessageSource.EXPECT().ListenMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, c chan *packet.Message) error {
		return nil
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "publish",
			fields: fields{
				ctx:           ctx,
				client:        NewClient(mockWriter, metaTopic),
				meta:          metaTopic,
				messageSource: mockMessageSource,
				cancel:        cancel,
				messageChan:   listenChan,
			},
			args: args{
				publish: &packet.Message{
					PublishPacket: &publishPacket,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS0{
				ctx:           tt.fields.ctx,
				cancel:        tt.fields.cancel,
				topic:         tt.fields.topic,
				messageSource: tt.fields.messageSource,
				client:        tt.fields.client,
				meta:          tt.fields.meta,
				messageChan:   tt.fields.messageChan,
			}
			go func() {
				_ = q.Start(tt.fields.ctx)
			}()
			time.Sleep(1 * time.Second)
			if err := q.Publish(tt.args.publish); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQoS0_Close(t *testing.T) {
	var (
		ctxClosed, cancel1 = context.WithCancel(context.TODO())
		ctx, cancel        = context.WithCancel(context.TODO())
	)
	cancel1()
	type fields struct {
		ctx           context.Context
		cancel        context.CancelFunc
		topic         string
		messageSource broker.MessageSource
		client        *Client
		meta          *topic.Meta
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "close",
			fields: fields{
				ctx:           ctxClosed,
				cancel:        nil,
				topic:         "",
				messageSource: nil,
				client:        nil,
				meta:          nil,
			},
			wantErr: true,
		},
		{
			name: "normal close",
			fields: fields{
				ctx:    ctx,
				cancel: cancel,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS0{
				ctx:           tt.fields.ctx,
				cancel:        tt.fields.cancel,
				topic:         tt.fields.topic,
				messageSource: tt.fields.messageSource,
				client:        tt.fields.client,
				meta:          tt.fields.meta,
			}
			if err := q.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQoS0_Start(t *testing.T) {
	var (
		ctx, cancel      = context.WithTimeout(context.TODO(), 2*time.Second)
		mockCtrl         = gomock.NewController(t)
		messageChan      = make(chan *packet.Message, 10)
		mockPacketWriter = client2.NewMockPacketWriter(mockCtrl)
		topicMeta        = &topic.Meta{}
	)
	messageChan <- &packet.Message{
		PublishPacket: &packets.Publish{},
	}

	defer func() {
		mockCtrl.Finish()
		cancel()
	}()

	mockMessageSource := broker2.NewMockMessageSource(mockCtrl)
	mockMessageSource.EXPECT().ListenMessage(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, c chan *packet.Message) error {
		return nil
	}).AnyTimes()
	mockPacketWriter.EXPECT().NextPacketID().Return(uint16(1)).AnyTimes()
	mockPacketWriter.EXPECT().WritePacket(gomock.Any()).Return(nil).AnyTimes()
	mockPacketWriter.EXPECT().GetID().Return("1").AnyTimes()

	type fields struct {
		ctx           context.Context
		cancel        context.CancelFunc
		topic         string
		messageSource broker.MessageSource
		client        *Client
		meta          *topic.Meta
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "start",
			fields: fields{
				topic:         "",
				messageSource: mockMessageSource,
				client:        NewClient(mockPacketWriter, topicMeta),
				meta:          topicMeta,
			},
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS0{
				ctx:           tt.fields.ctx,
				cancel:        tt.fields.cancel,
				topic:         tt.fields.topic,
				messageSource: tt.fields.messageSource,
				client:        tt.fields.client,
				meta:          tt.fields.meta,
			}
			if err := q.Start(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
