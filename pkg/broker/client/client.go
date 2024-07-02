package client

import (
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
)

type ID interface {
	GetID() string
}

type PacketIDGenerator interface {
	NextPacketID() uint16
}

type WritePacket struct {
	packets.Packet
	// FullTopic is the full topic of the message.
	// If the message has a shared topic.
	// The FullTopic is the shared topic.
	// Like $share/shareGroup/topic
	FullTopic string
}

func NewWritePacket(p packets.Packet) *WritePacket {
	return &WritePacket{
		Packet: p,
	}
}

type PacketWriter interface {
	// WritePacket writes the packet to the writer.
	// Warning: packetID is original packetID, method should change it to the new one that does not used.
	WritePacket(packet *WritePacket) (err error)
	PacketIDGenerator

	ID
	Close() error
}

type Client interface {
	Publish(publish *packet.Message) error
	PubRel(message *packet.Message) error

	GetUnFinishedMessage() []*packet.Message
	GetPacketWriter() PacketWriter
	HandlePublishResponse
}

type HandlePublishResponse interface {
	HandlePublishAck(pubAck *packets.Puback) (ok bool, err error)
	HandlePublishRec(pubRec *packets.Pubrec) (ok bool, err error)

	// HandelPublishComp HandlePublishComp handles the PUBCOMP packet from the client.
	// ok is true if the packet exists.
	// err is the error that occurred.
	HandelPublishComp(pubComp *packets.Pubcomp) (ok bool, err error)
}
