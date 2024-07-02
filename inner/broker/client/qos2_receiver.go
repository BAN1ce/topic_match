package client

import (
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"go.uber.org/zap"
	"sync"
)

// QoS2ReceiveStore is a store for one client's QoS2 packets
type QoS2ReceiveStore struct {
	mux     sync.Mutex
	waiting map[uint16]*packet.Message
}

func NewQoS2Handler() *QoS2ReceiveStore {
	return &QoS2ReceiveStore{
		waiting: make(map[uint16]*packet.Message),
	}
}

// StoreNotExists handles the publishing QoS2 packet from the client
func (w *QoS2ReceiveStore) StoreNotExists(publish *packets.Publish) bool {
	w.mux.Lock()
	defer w.mux.Unlock()

	if _, ok := w.waiting[publish.PacketID]; ok {
		return false
	}

	w.waiting[publish.PacketID] = &packet.Message{
		PublishPacket: publish,
	}
	logger.Logger.Debug("QoS2ReceiveStore: handle publish", zap.Uint16("packet id", publish.PacketID))
	return true
}

// Delete HandlePubRec handles the PubRec packet from the client
// If the packet UID is not found, it will return false
func (w *QoS2ReceiveStore) Delete(pubRel *packets.Pubrel) (*packets.Publish, bool) {
	w.mux.Lock()
	defer w.mux.Unlock()

	if message, ok := w.waiting[pubRel.PacketID]; !ok {
		return nil, false
	} else {
		delete(w.waiting, pubRel.PacketID)
		return message.PublishPacket, true
	}
}
