package sub_topic

import (
	"context"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	topic2 "github.com/BAN1ce/skyTree/pkg/mocks/topic"
	"github.com/eclipse/paho.golang/packets"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestNewManager(t *testing.T) {
	var (
		result = &Topics{
			topic: map[string]topic.Topic{},
		}
	)
	type args struct {
		ctx context.Context
		ops []Option
	}
	tests := []struct {
		name string
		args args
		want *Topics
	}{
		{
			name: "new manager success",
			args: args{
				ctx: context.TODO(),
				ops: nil,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTopics(tt.args.ctx, tt.args.ops...); !reflect.DeepEqual(got.topic, tt.want.topic) {
				t.Errorf("NewTopics() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_HandlePublishAck(t1 *testing.T) {
	var (
		mockCtrl     = gomock.NewController(t1)
		mockQosTopic = topic2.NewMockQoS1Subscriber(mockCtrl)
		topics       = map[string]topic.Topic{
			"mock": mockQosTopic,
		}
	)

	mockQosTopic.EXPECT().HandlePublishAck(gomock.Any()).Return(true, nil).AnyTimes()

	type fields struct {
		ctx    context.Context
		cancel context.CancelCauseFunc
		topic  map[string]topic.Topic
	}
	type args struct {
		topicName string
		puback    *packets.Puback
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "not exists",
			fields: fields{
				ctx:    nil,
				cancel: nil,
				topic:  topics,
			},
			args: args{
				topicName: uuid.NewString(),
				puback:    nil,
			},
			wantErr: true,
		},
		{
			name: "handle publish ack",
			fields: fields{
				ctx:    nil,
				cancel: nil,
				topic:  topics,
			},
			args: args{
				topicName: "mock",
				puback: &packets.Puback{
					PacketID: 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Topics{
				ctx:    tt.fields.ctx,
				cancel: tt.fields.cancel,
				topic:  tt.fields.topic,
			}
			if _, err := t.HandlePublishAck(tt.args.topicName, tt.args.puback); (err != nil) != tt.wantErr {
				t1.Errorf("HandlePublishAck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestManager_DeleteTopic(t1 *testing.T) {
	var (
		mockCtrl = gomock.NewController(t1)
		mockQoS  = topic2.NewMockQoS1Subscriber(mockCtrl)
		closed   bool
		topics   = map[string]topic.Topic{
			"mock": mockQoS,
		}
	)
	defer mockCtrl.Finish()

	mockQoS.EXPECT().Close().DoAndReturn(func() error {
		closed = true
		return nil

	}).AnyTimes()

	type fields struct {
		ctx    context.Context
		cancel context.CancelCauseFunc
		topic  map[string]topic.Topic
	}
	type args struct {
		topicName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "delete topic",
			fields: fields{
				ctx:    nil,
				cancel: nil,
				topic:  topics,
			},
			args: args{
				"mock",
			},
		},
		{
			name: "delete topic, not found",
			fields: fields{
				ctx:    nil,
				cancel: nil,
				topic:  topics,
			},
			args: args{
				uuid.NewString(),
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Topics{
				ctx:    tt.fields.ctx,
				cancel: tt.fields.cancel,
				topic:  tt.fields.topic,
			}
			t.DeleteTopic(tt.args.topicName)
			if !closed {
				t1.Errorf("DeleteTopic() error, didn't close")
			}
		})
	}
}
