package packet

import (
	"github.com/BAN1ce/skyTree/config"
	"github.com/eclipse/paho.golang/packets"
	"sync/atomic"
	"time"
)

const (
	PubReceived = 1 << iota
	FromSession
	Will
)

type Message struct {
	SendClientID string
	MessageID    string
	PacketID     string
	// PublishPacket should never change after it is set
	PublishPacket *packets.Publish `json:"-"`

	// PubReceived is used to determine whether the message has been received by the client
	// Sender -> Receiver: Publish
	// Receiver -> Sender: PubRec
	PubReceived bool
	// CreatedTime is the time when the message is created, unit is nanosecond
	CreatedTime    int64
	ExpiredTime    int64
	Will           bool
	State          int
	ShareTopic     string
	SubscribeTopic string
	OwnerClientID  string
	Duplicate      bool
	Retain         bool
	RetryInfo      *RetryInfo
}

func (m *Message) IsExpired() bool {
	return time.Now().Unix() > m.ExpiredTime
}

type RetryInfo struct {
	Key          string
	Times        atomic.Int32
	FirstPubTime time.Time
	IntervalTime time.Duration
}

func (r *RetryInfo) IsTimeout() bool {
	cfg := config.GetConfig()
	if int(r.Times.Load()) >= cfg.Retry.GetMaxTimes() {
		return true
	}

	if time.Since(r.FirstPubTime) > cfg.Retry.GetTimeout() {
		return true
	}
	return false
}

type MessageExtraInfo struct {
	Duplicate bool
	Retain    bool
}

func (m *Message) DeepCopy() *Message {
	return &Message{
		SendClientID:   m.SendClientID,
		MessageID:      m.MessageID,
		PacketID:       m.PacketID,
		PublishPacket:  m.PublishPacket,
		PubReceived:    m.PubReceived,
		CreatedTime:    m.CreatedTime,
		ExpiredTime:    m.ExpiredTime,
		Will:           m.Will,
		State:          m.State,
		ShareTopic:     m.ShareTopic,
		SubscribeTopic: m.SubscribeTopic,
		OwnerClientID:  m.OwnerClientID,
	}
}

func (m *Message) SetSubIdentifier(id byte) {
	if m.PublishPacket.Properties == nil {
		m.PublishPacket.Properties = &packets.Properties{}
	}
	m.PublishPacket.Properties.SubIDAvailable = &id
}

func (m *Message) SetPubReceived(received bool) {
	if received {
		m.State |= PubReceived
	} else {
		m.State &= ^PubReceived
	}
}

func (m *Message) SetFromSession(fromSession bool) {
	if fromSession {
		m.State |= FromSession
	} else {
		m.State &= ^FromSession
	}
}

func (m *Message) SetWill(will bool) {
	if will {
		m.State |= Will
	} else {
		m.State &= ^Will
	}
}

//

func (m *Message) IsPubReceived() bool {
	return m.State&PubReceived == PubReceived
}

func (m *Message) IsFromSession() bool {
	return m.State&FromSession == FromSession
}

func (m *Message) IsWill() bool {
	return m.State&Will == Will
}

func (m *Message) ToSessionPayload() string {

	return ""
}

func (m *Message) Release() {
	// TODO: release
}

func (m *Message) GetFullTopic() string {
	if m.ShareTopic != "" {
		return m.ShareTopic
	}
	return m.PublishPacket.Topic
}
