package pool

import (
	packet2 "github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"sync"
)

var (
	PublRel = NewPubRel()
)

type PubRel struct {
	sync.Pool
}

func (p *PubRel) Get() *packets.Pubrel {
	return p.Pool.Get().(*packets.Pubrel)
}

func (p *PubRel) Put(b *packets.Pubrel) {
	b.Properties = nil
	p.Pool.Put(b)
}

func NewPubRel() *PubRel {
	return &PubRel{
		sync.Pool{
			New: func() interface{} {
				return packet2.NewPublishRel()
			},
		},
	}
}
