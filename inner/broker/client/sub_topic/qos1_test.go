package sub_topic

import (
	"context"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	client2 "github.com/BAN1ce/skyTree/pkg/mocks/client"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestQoS1_Meta(t *testing.T) {
	var (
		metaData = topic.Meta{
			Topic: "test",
		}
	)
	type fields struct {
		ctx               context.Context
		cancel            context.CancelFunc
		meta              topic.Meta
		publishChan       *PublishChan
		client            *WithRetryClient
		messageSource     broker.MessageSource
		unfinishedMessage []*packet.Message
	}
	tests := []struct {
		name   string
		fields fields
		want   topic.Meta
	}{
		{
			name: "meta",
			fields: fields{
				meta: metaData,
			},
			want: metaData,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS1{
				ctx:               tt.fields.ctx,
				cancel:            tt.fields.cancel,
				meta:              tt.fields.meta,
				publishChan:       tt.fields.publishChan,
				client:            tt.fields.client,
				messageSource:     tt.fields.messageSource,
				unfinishedMessage: tt.fields.unfinishedMessage,
			}
			if got := q.Meta(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Meta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQoS1_GetUnFinishedMessage(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockClient = client2.NewMockClient(mockCtrl)
		unMessage  = []*packet.Message{
			{},
		}
	)
	defer mockCtrl.Finish()

	mockClient.EXPECT().GetUnFinishedMessage().Return(unMessage).AnyTimes()

	type fields struct {
		client client.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   []*packet.Message
	}{
		{
			name: "GetUnFinishedMessage",
			fields: fields{
				client: mockClient,
			},
			want: unMessage,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS1{
				client: tt.fields.client,
			}
			if got := q.GetUnFinishedMessage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUnFinishedMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQoS1_Publish(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.TODO())
		pubValue    = &packet.Message{}
	)
	defer cancel()
	type fields struct {
		ctx         context.Context
		publishChan *PublishChan
	}
	type args struct {
		publish *packet.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "publish",
			fields: fields{
				ctx:         ctx,
				publishChan: NewPublishChan(10),
			},
			args: args{
				pubValue,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &QoS1{
				ctx:         tt.fields.ctx,
				publishChan: tt.fields.publishChan,
			}
			if err := q.Publish(tt.args.publish); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(pubValue, <-q.publishChan.ch) {
				t.Errorf("not match")
			}
		})
	}
}
