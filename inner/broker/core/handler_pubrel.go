package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/eclipse/paho.golang/packets"
)

type PublishRel struct {
}

func NewPublishRel() *PublishRel {
	return &PublishRel{}
}

func (p *PublishRel) Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error) {
	return nil
}
