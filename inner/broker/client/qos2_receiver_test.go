package client

import (
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"reflect"
	"testing"
)

func TestNewQoS2Handler(t *testing.T) {
	tests := []struct {
		name string
		want *QoS2ReceiveStore
	}{
		{
			name: "new",
			want: &QoS2ReceiveStore{
				waiting: make(map[uint16]*packet.Message),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewQoS2Handler(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQoS2Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQoS2ReceiveStore_StoreNotExists(t *testing.T) {
	var (
		packetID uint16 = 1
		waiting         = map[uint16]*packet.Message{
			2: {},
		}
	)
	type fields struct {
		waiting map[uint16]*packet.Message
	}
	type args struct {
		publish *packets.Publish
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "store not exists",
			fields: fields{
				waiting: waiting,
			},
			args: args{
				publish: &packets.Publish{
					PacketID: packetID,
				},
			},
			want: true,
		},
		{
			name: "store exists",
			fields: fields{
				waiting: waiting,
			},
			args: args{
				publish: &packets.Publish{
					PacketID: 2,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &QoS2ReceiveStore{
				waiting: tt.fields.waiting,
			}
			if got := w.StoreNotExists(tt.args.publish); got != tt.want {
				t.Errorf("StoreNotExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQoS2ReceiveStore_Delete(t *testing.T) {
	var (
		packetID uint16 = 1
		message         = &packet.Message{
			PublishPacket: &packets.Publish{
				PacketID: packetID,
			},
		}
		waiting = map[uint16]*packet.Message{
			packetID: message,
		}
	)

	type fields struct {
		waiting map[uint16]*packet.Message
	}
	type args struct {
		pubRel *packets.Pubrel
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *packets.Publish
		want1  bool
	}{
		{
			name: "delete",
			fields: fields{
				waiting: waiting,
			},
			args: args{
				pubRel: &packets.Pubrel{
					PacketID: packetID,
				},
			},
			want:  message.PublishPacket,
			want1: true,
		},

		{
			name: "delete not exists",
			fields: fields{
				waiting: waiting,
			},
			args: args{
				pubRel: &packets.Pubrel{
					PacketID: packetID + 1,
				},
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &QoS2ReceiveStore{
				waiting: tt.fields.waiting,
			}
			got, got1 := w.Delete(tt.args.pubRel)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Delete() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
