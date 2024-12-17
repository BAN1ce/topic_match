package middleware

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	client2 "github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/eclipse/paho.golang/packets"
)

type OnceConnect struct {
}

func (o *OnceConnect) Handle(c *client.Client, packet *packets.ControlPacket) error {
	if c.IsState(client.ReceivedConnect) {
		disconnect := packets.NewControlPacket(packets.DISCONNECT).Content.(*packets.Disconnect)
		disconnect.ReasonCode = packets.ConnackProtocolError
		c.Write(&client2.WritePacket{
			Packet: disconnect,
		})
		return errs.ErrConnectPacketDuplicate
	}
	return nil
}
