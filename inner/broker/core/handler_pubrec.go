package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
	"github.com/eclipse/paho.golang/packets"
)

type PublishRec struct {
}

func NewPublishRec() *PublishRec {
	return &PublishRec{}
}

func (p *PublishRec) Handle(broker *Broker, client *client.Client, rawPacket *packets.ControlPacket) (err error) {
	return err
}
