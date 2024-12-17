package store

import packet2 "github.com/BAN1ce/skyTree/pkg/packet"

type Wrapper interface {
	StorePublishPacket(topics map[string]int32, packet *packet2.Message) (messageID string, err error)
}
