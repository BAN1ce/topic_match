package model

import (
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
)

type Session struct {
	ID              string                     `json:"id"`
	WillMessage     *session.WillMessage       `json:"will_message"`
	ConnectProperty *session.ConnectProperties `json:"connect_property"`
	ExpiryInterval  int64                      `json:"expiry_interval"`
	Topic           []*topic.Meta              `json:"topic"`
}

func ToSessionModel(session2 session.Session) *Session {
	var (
		result = &Session{}
	)
	result.ExpiryInterval = session2.GetExpiryInterval()
	if properties, err := session2.GetConnectProperties(); err == nil {
		result.ConnectProperty = properties
	}
	if will, ok, err := session2.GetWillMessage(); ok && err == nil {
		result.WillMessage = will
	}
	topics := session2.ReadSubTopics()
	result.Topic = topics

	return result

}
