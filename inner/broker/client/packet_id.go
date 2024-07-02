package client

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"sync"
)

type PacketIDFactory struct {
	id  uint16
	mux sync.RWMutex
}

func NewPacketIDFactory() *PacketIDFactory {
	var (
		id = randomInitPacketID()
	)
	return &PacketIDFactory{
		id: id,
	}
}

func randomInitPacketID() uint16 {
	return uint16(rand.Intn(math.MaxUint16))
}

func (p *PacketIDFactory) SetID(id uint16) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.id = id
}

func (p *PacketIDFactory) NextPacketID() uint16 {
	p.mux.Lock()
	defer p.mux.Unlock()

	if p.id == math.MaxUint16 {
		p.id = 0
	} else {
		p.id++
	}
	return p.id
}

func PacketIDToString(id uint16) string {
	return fmt.Sprintf("%d", id)
}

func StringToPacketID(id string) uint16 {
	i, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0
	}
	return uint16(i)
}

type PacketIDTopic struct {
	mux      sync.RWMutex
	packetID map[uint16]string
}

func NewPacketIDTopic() *PacketIDTopic {
	return &PacketIDTopic{
		packetID: make(map[uint16]string),
	}
}

func (p *PacketIDTopic) SePacketIDTopic(id uint16, topic string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.packetID[id] = topic
}

func (p *PacketIDTopic) GetTopic(id uint16) string {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.packetID[id]
}

func (p *PacketIDTopic) DeletePacketID(id uint16) {
	p.mux.Lock()
	defer p.mux.Unlock()
	delete(p.packetID, id)
}
