package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/eclipse/paho.golang/packets"
)

type PingHandler struct {
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (p *PingHandler) Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error) {
	return nil
}
