package sub_topic

import (
	"context"
	"github.com/BAN1ce/skyTree/pkg/broker"
	broker2 "github.com/BAN1ce/skyTree/pkg/mocks/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func TestFillUnfinishedMessage(t *testing.T) {
	var (
		mockCtrl          = gomock.NewController(t)
		mockMessageSource = broker2.NewMockMessageSource(mockCtrl)
		publishTest       = &packets.Publish{
			Topic: "test",
		}
	)
	defer mockCtrl.Finish()
	mockMessageSource.EXPECT().NextMessages(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]*packet.Message{
		{
			MessageID:     "1",
			PublishPacket: publishTest,
		},
	}, nil).AnyTimes()

	type args struct {
		ctx     context.Context
		message []*packet.Message
		source  broker.MessageSource
	}
	tests := []struct {
		name string
		args args
		want []*packet.Message
	}{
		{
			name: "Test FillUnfinishedMessage",
			args: args{
				ctx: context.Background(),
				message: []*packet.Message{
					{
						MessageID: "1",
					},
				},
				source: mockMessageSource,
			},
			want: []*packet.Message{
				{
					MessageID:     "1",
					PublishPacket: publishTest,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FillUnfinishedMessage(tt.args.ctx, tt.args.message, tt.args.source); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillUnfinishedMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
