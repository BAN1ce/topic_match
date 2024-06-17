package client

import (
	"github.com/BAN1ce/skyTree/logger"
	"github.com/eclipse/paho.golang/packets"
	"go.uber.org/zap"
	"sync"
)

type QoS2Message struct {
	publishPacket *packets.Publish
}

type QoS2Handler struct {
	mux     sync.Mutex
	waiting map[uint16]*QoS2Message
}

func NewQoS2Handler() *QoS2Handler {
	return &QoS2Handler{
		waiting: make(map[uint16]*QoS2Message),
	}
}

func (w *QoS2Handler) HandlePublish(publish *packets.Publish) bool {
	w.mux.Lock()
	defer w.mux.Unlock()

	if _, ok := w.waiting[publish.PacketID]; ok {
		return false
	}

	w.waiting[publish.PacketID] = &QoS2Message{
		publishPacket: publish,
	}
	logger.Logger.Debug("QoS2Handler: handle publish", zap.Uint16("packet id", publish.PacketID))
	return true
}

func (w *QoS2Handler) HandlePubRel(pubRel *packets.Pubrel) (*packets.Publish, bool) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if message, ok := w.waiting[pubRel.PacketID]; !ok {
		return nil, false
	} else {
		delete(w.waiting, pubRel.PacketID)
		return message.publishPacket, true
	}
}
