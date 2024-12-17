package utils

import (
	"github.com/BAN1ce/Tree/proto"
	"github.com/eclipse/paho.golang/packets"
	"github.com/google/uuid"
	"math/rand"
	"strings"
	"time"
)

func SplitTopic(topic string) []string {
	if len(topic) == 0 {
		return nil
	}
	var result = make([]string, 0, 30)

	if topic[0] == '/' {
		result = append(result, "/")

	}
	//if !strings.Contains(topic, "/") {
	//	result = append(result, "/")
	//}

	tmp := strings.Split(strings.Trim(topic, "/"), "/")

	for _, v := range tmp {
		result = append(result, v)
	}

	for _, v := range result {
		if v != "/" {
			return result
		}
	}

	return nil

}

// nolint
func HasWildcard(topic string) bool {
	return strings.Contains(topic, "+") || strings.Contains(topic, "#")
}

// nolint
func subQosMoreThan0(topics map[string]int32) bool {
	for _, v := range topics {
		if v > 0 {
			return true
		}
	}
	return false
}

func ParseShareTopic(shareTopic string) (shareGroup, subTopic string) {
	shareNameSubTopic := strings.TrimPrefix(shareTopic, "$share/")
	index := strings.Index(shareNameSubTopic, "/")
	if index == -1 {
		return "", ""
	}
	return shareNameSubTopic[:index], shareNameSubTopic[index+1:]
}

func IsShareTopic(shareTopic string) bool {
	return strings.HasPrefix(shareTopic, "$share/")
}

func SubOptionToProtoSubOption(options *packets.SubOptions) *proto.SubOption {
	return &proto.SubOption{
		QoS:               ByteToInt32(options.QoS),
		NoLocal:           options.NoLocal,
		RetainAsPublished: options.RetainAsPublished,
	}
}

func RandomWildcardTopic() string {
	var (
		result       = "/"
		randomLength = rand.Intn(20)
	)

	for i := 0; i < randomLength; i++ {
		switch time.Now().Nanosecond() % 3 {
		case 0:
			result += "/+"
		case 1:
			result += "/#"
			return result
		case 2:
			result += "/" + uuid.NewString()

		}

	}
	return result
}
