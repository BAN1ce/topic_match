package sub_topic

import (
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
	"github.com/eclipse/paho.golang/packets"
)

type Client struct {
	writer client.PacketWriter
	meta   *topic.Meta
}

func NewClient(writer client.PacketWriter, meta *topic.Meta) *Client {
	return &Client{
		writer: writer,
		meta:   meta,
	}
}

func (c *Client) addSubTopicToProperties(pub *packets.Publish, subTopic string) {
	if pub.Properties == nil {
		pub.Properties = &packets.Properties{}
	}
	pub.Properties.User = []packets.User{
		{
			Key:   "sub_topic",
			Value: subTopic,
		},
	}
}

// Publish publishes the message to the client.
// And should not change the publishing packet
func (c *Client) Publish(publish *packet.Message) error {
	// NoLocal means that the server does not send messages published by the client itself.
	if c.meta.NoLocal && c.writer.GetID() == publish.SendClientID {
		return nil
	}

	var (
		newPublishPacket = pool.PublishPool.Get()
	)

	defer pool.PublishPool.Put(newPublishPacket)

	pool.CopyPublish(newPublishPacket, publish.PublishPacket)

	// if the topic is different from the subscription topic, add the subscription topic to the properties
	// example: subscription topic is "a/b/#", the topic of the message is "a/b/c", then add "a/b/#" to the properties
	if publish.PublishPacket.Topic != publish.SubscribeTopic {
		c.addSubTopicToProperties(newPublishPacket, publish.SubscribeTopic)
	}

	// if a topic has an identifier, set the identifier to the publishing packet
	if c.meta.Identifier != 0 {
		identifier := byte(c.meta.Identifier)
		if newPublishPacket.Properties == nil {
			newPublishPacket.Properties = &packets.Properties{}
		}
		newPublishPacket.Properties.SubIDAvailable = &identifier
	}
	var (
		fullTopic = newPublishPacket.Topic
	)
	if publish.ShareTopic != "" {
		fullTopic = publish.ShareTopic
	}

	if publish.SubscribeTopic != "" {
		fullTopic = publish.SubscribeTopic
	}

	//if publish.Duplicate {
	//	return c.writer.RetryWrite(&client.WritePacket{
	//		Packet:    publish.PublishPacket,
	//		FullTopic: fullTopic,
	//	})
	//}

	return c.writer.Write(&client.WritePacket{
		Packet:    newPublishPacket,
		FullTopic: fullTopic,
	})
}

// PubRel Sends a PUBREL packet to the client.
// For response to a PUBREC packet.
func (c *Client) PubRel(message *packet.Message) error {
	pubRel := pool.PublRel.Get()
	defer pool.PublRel.Put(pubRel)

	pubRel.PacketID = message.PublishPacket.PacketID
	return c.writer.Write(&client.WritePacket{
		Packet: pubRel,
	})
}

func (c *Client) GetPacketWriter() client.PacketWriter {
	return c.writer
}

func (c *Client) HandlePublishAck(pubAck *packets.Puback) (err error) {
	// do nothing
	return err
}

func (c *Client) HandlePublishRec(pubRec *packets.Pubrec) (err error) {
	// do nothing
	return err
}

func (c *Client) HandelPublishComp(pubComp *packets.Pubcomp) (err error) {
	// do nothing
	return err
}

func (c *Client) GetUnFinishedMessage() []*packet.Message {
	return nil
}

func (c *Client) HandlePubRel(pubrel *packets.Pubrel) {
	//do nothing
}
