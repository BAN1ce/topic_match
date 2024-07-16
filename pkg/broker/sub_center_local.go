package broker

import (
	"github.com/BAN1ce/Tree/inner/api"
	"github.com/BAN1ce/Tree/inner/store"
	"github.com/BAN1ce/Tree/proto"
	"github.com/BAN1ce/skyTree/pkg/utils"
	"github.com/eclipse/paho.golang/packets"
)

type LocalSubCenter struct {
	state *store.State
}

func NewLocalSubCenter() *LocalSubCenter {
	return &LocalSubCenter{
		state: store.NewState(),
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

	_, err := l.state.HandleSubRequest(&proto.SubRequest{
		Topics:   topicsRequest,
		ClientID: clientID,
	})
	return err
}

func (l *LocalSubCenter) DeleteSub(clientID string, topics []string) error {
	_, err := l.state.HandleUnSubRequest(&proto.UnSubRequest{
		Topics:   topics,
		ClientID: clientID,
	})
	return err
}

func (l *LocalSubCenter) Match(topic string) (clientIDQos map[string]int32) {
	rsp := l.state.MatchTopic(&api.MatchTopicRequest{
		Topic: topic,
	})

	clientIDQos = map[string]int32{}
	for _, client := range rsp.Client {
		clientIDQos[client.ClientID] = client.QoS
	}
	return
}

func (l *LocalSubCenter) MatchTopic(topic string) (topics map[string]int32) {
	rsp := l.state.MatchSubTopics(&api.MatchSubTopicRequest{
		Topic: topic,
	})
	return rsp.Topic
}

func (l *LocalSubCenter) DeleteClient(clientID string) error {
	//TODO implement me
	panic("implement me")
}
