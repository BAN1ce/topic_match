package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/eclipse/paho.golang/packets"
)

type UnsubHandler struct {
}

func NewUnsubHandler() *UnsubHandler {
	return &UnsubHandler{}
}
func (u *UnsubHandler) Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error) {
	return err
}
