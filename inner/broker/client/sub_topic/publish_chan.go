package sub_topic

import (
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/packet"
)

type PublishChan struct {
	ch chan *packet.Message
}

func NewPublishChan(size int) *PublishChan {
	return &PublishChan{
		ch: make(chan *packet.Message, size),
	}

}
func (q *PublishChan) writeToPublishWithDrop(message *packet.Message) {
	select {
	case q.ch <- message:

	default:
		// TODO: add metric for dropped message
		logger.Logger.Warn().Msg("QoS1: ch is full, message will be dropped")

	}
}

func (q *PublishChan) write(message *packet.Message) {
	q.ch <- message
}
