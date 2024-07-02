package sub_topic

import (
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	client2 "github.com/BAN1ce/skyTree/pkg/mocks/client"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestClient_PubRel(t *testing.T) {
	var (
		mockCtrl   = gomock.NewController(t)
		mockWriter = client2.NewMockPacketWriter(mockCtrl)
		pubRel     = &packets.Pubrel{
			PacketID: 1,
		}
	)
	defer mockCtrl.Finish()

	mockWriter.EXPECT().WritePacket(gomock.Eq(&client.WritePacket{
		Packet: pubRel,
	})).Return(nil).AnyTimes()
	type fields struct {
		writer client.PacketWriter
		meta   *topic.Meta
	}
	type args struct {
		message *packet.Message
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "pubrel success",
			fields: fields{
				writer: mockWriter,
				meta:   nil,
			},
			args: args{
				message: &packet.Message{
					PubRelPacket: pubRel,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				writer: tt.fields.writer,
				meta:   tt.fields.meta,
			}
			if err := c.PubRel(tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("PubRel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
