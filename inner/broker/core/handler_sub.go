package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/eclipse/paho.golang/packets"
)

type SubHandler struct {
}

func NewSubHandler() *SubHandler {
	return &SubHandler{}
}

func (s *SubHandler) Handle(b *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error) {
	return err
}
