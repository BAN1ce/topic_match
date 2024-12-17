package broker

import (
	"context"
	"github.com/BAN1ce/Tree/app"
	"github.com/BAN1ce/Tree/inner/api/request"
	"github.com/BAN1ce/Tree/inner/core"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"github.com/eclipse/paho.golang/packets"
)

type LocalSubCenter struct {
	state *core.Core
}

func NewLocalSubCenter() *LocalSubCenter {
	newApp := app.NewApp()
	return &LocalSubCenter{
		state: newApp.State,
	}
}
func (l *LocalSubCenter) CreateSub(clientID string, topics []packets.SubOptions) error {
	var topicsRequest = make(map[string]*proto.SubOption)
	for _, opt := range topics {
		var (
			shareTopic string
			subOption  = &proto.SubOption{
				QoS:               int32(opt.QoS),
				NoLocal:           opt.NoLocal,
				RetainAsPublished: opt.RetainAsPublished,
				Topic:             opt.Topic,
			}
		)
		if utils.IsShareTopic(opt.Topic) {
			subOption.ShareGroup, shareTopic = utils.ParseShareTopic(opt.Topic)
			subOption.Share = true
			subOption.Topic = shareTopic
			topicsRequest[shareTopic] = subOption
		}
		topicsRequest[opt.Topic] = subOption
	}

	l.state.HandleSubRequest(&proto.SubRequest{
		Topics:   topicsRequest,
		ClientID: clientID,
	})
	return nil
}

func (l *LocalSubCenter) DeleteSub(clientID string, topics []string) error {
	l.state.HandleUnSubRequest(&proto.UnSubRequest{
		Topics:   topics,
		ClientID: clientID,
	})
	return nil
}

func (l *LocalSubCenter) Match(topic string) (clientIDQos map[string]int32) {
	panic("implement me")
}

func (l *LocalSubCenter) GetAllMatchTopics(topic string) (topics map[string]int32) {
	rsp := l.state.GetAllMatchTopics(context.TODO(), &request.GetAllMatchTopicsRequest{
		Topic: topic,
	})
	return rsp.Topic
}

func (l *LocalSubCenter) DeleteClient(clientID string) error {
	//TODO implement me
	panic("implement me")
}

func (l *LocalSubCenter) GetMatchTopicForWildcardTopic(wildTopic string) []string {
	response := l.state.GetMatchTopicForWildcardTopic(context.TODO(), &request.GetAllMatchTopicsForWildTopicRequest{
		Topic: wildTopic,
	})
	return response.Topic
}
