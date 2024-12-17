package topic

import (
	"context"
	"encoding/json"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"github.com/eclipse/paho.golang/packets"
	"strings"
	"sync"
)

type Topic interface {
	Start(ctx context.Context) error
	Close() error
	Publish(publish *packet.Message, extra *packet.MessageExtraInfo) error
	GetUnFinishedMessage() []*packet.Message
	Meta() *Meta
}

type QoS0Subscriber interface {
	Topic
}

type QoS1Subscriber interface {
	Topic
	HandlePublishAck(puback *packets.Puback) (err error)
}

type QoS2Subscriber interface {
	Topic
	HandlePublishRec(pubrec *packets.Pubrec) (err error)
	HandlePublishComp(pubcomp *packets.Pubcomp) (err error)
}

type Meta struct {
	mux               sync.RWMutex
	Identifier        int            `json:"identifier,omitempty"`
	Topic             string         `json:"topic"`
	NoLocal           bool           `json:"no_local"`
	RetainAsPublished bool           `json:"retain_as_published"`
	RetainHandling    int            `json:"retain_handling"`
	WindowSize        int            `json:"window_size"`
	LatestMessageID   string         `json:"latest_message_id"`
	QoS               int32          `json:"qos"`
	Properties        []packets.User `json:"properties"`
	Share             bool           `json:"share"`
	ShareTopic        *ShareTopic    `json:"share_topic"`
}

func (m *Meta) SetLatestMessageID(id string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.LatestMessageID = id
}

func (m *Meta) Meta() *Meta {
	m.mux.RLock()
	defer m.mux.RUnlock()
	var tmp Meta

	b, _ := json.Marshal(m)
	json.Unmarshal(b, &tmp)
	return &tmp
}

func (m *Meta) String() string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	b, _ := json.Marshal(m)
	return string(b)
}

func NewMetaFromSubPacket(subOption *packets.SubOptions, properties *packets.Properties) *Meta {
	m := &Meta{
		Topic:             subOption.Topic,
		NoLocal:           subOption.NoLocal,
		RetainAsPublished: subOption.RetainAsPublished,
		RetainHandling:    int(subOption.RetainHandling),
		QoS:               int32(subOption.QoS),
		Properties:        properties.User,
	}
	if utils.IsShareTopic(subOption.Topic) {
		m.Share = true
		m.ShareTopic = NewShareTopic(subOption.Topic)
	}

	if properties.SubscriptionIdentifier != nil {
		m.Identifier = *properties.SubscriptionIdentifier
	}
	return m
}

func SplitShareAndNoShare(subPacket *packets.Subscribe) (shareSubscribe *packets.Subscribe, noShareSubscribe *packets.Subscribe) {
	var (
		shareTopic   = make([]packets.SubOptions, 0)
		noShareTopic = make([]packets.SubOptions, 0)
	)
	shareSubscribe = &packets.Subscribe{
		PacketID:   subPacket.PacketID,
		Properties: subPacket.Properties,
	}
	noShareSubscribe = &packets.Subscribe{
		PacketID:   subPacket.PacketID,
		Properties: subPacket.Properties,
	}
	for _, topic := range subPacket.Subscriptions {
		if utils.IsShareTopic(topic.Topic) {
			shareTopic = append(shareTopic, topic)
		} else {
			noShareTopic = append(noShareTopic, topic)
		}
	}
	shareSubscribe.Subscriptions = shareTopic
	noShareSubscribe.Subscriptions = noShareTopic
	return shareSubscribe, noShareSubscribe
}

type ShareTopic struct {
	FullTopic  string
	ShareGroup string
	Topic      string
}

func NewShareTopic(fullTopic string) *ShareTopic {
	st := &ShareTopic{
		FullTopic: fullTopic,
	}
	s := strings.TrimPrefix(fullTopic, "$share/")

	if i := strings.Index(s, "/"); i != -1 {
		st.ShareGroup = s[0:i]
	}
	s = strings.TrimPrefix(s, st.ShareGroup)
	st.Topic = strings.TrimPrefix(s, "/")
	return st
}

func GetShareName(topic string) string {
	return strings.Split(topic, "/")[1]

}
func DeleteSharePrefix(fullTopic string) string {
	return strings.TrimPrefix(fullTopic, "$share")
}
