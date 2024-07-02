package utils

import (
	"github.com/BAN1ce/Tree/proto"
	"github.com/eclipse/paho.golang/packets"
	"strings"
)

func SplitTopic(topic string) []string {
	tmp := strings.Split(strings.Trim(topic, "/"), "/")
	for _, v := range tmp {
		if v == "" {
			result := make([]string, 0)
			for _, v := range tmp {
				if v != "" {
					result = append(result, v)
				}
			}
			return result
		}
	}
	return tmp
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
