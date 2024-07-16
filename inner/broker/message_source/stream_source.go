package message_source

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
	"sync"
)

type StreamSource struct {
	ctx             context.Context
	cancel          context.CancelFunc
	topic           string
	publishListener broker.PublishListener
	messageChan     chan *packet.Message
	mux             sync.RWMutex
}

// NewStreamSource creates a new event source with the given topic.
func NewStreamSource(topic string, listener broker.PublishListener) *StreamSource {
	return &StreamSource{
		topic:           topic,
		publishListener: listener,
	}
}

// Close closes the event source.
func (s *StreamSource) Close() error {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.ctx == nil {
		return fmt.Errorf("stream source doesn't listen to the topic %s", s.topic)
	}

	if s.ctx.Err() != nil {
		return nil
	}

	s.close()
	return nil
}

func (s *StreamSource) close() {
	s.cancel()
	s.publishListener.DeletePublishEvent(s.topic, s.handler)
	close(s.messageChan)

}
func (s *StreamSource) ListenMessage(ctx context.Context, writer chan *packet.Message) error {
	if s == nil {
		return fmt.Errorf("event source is nil")
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	s.messageChan = writer
	s.ctx, s.cancel = context.WithCancel(ctx)
	s.publishListener.CreateListenPublishEvent(s.topic, s.handler)
	return nil
}

// handler is the handler of the topic, it will be called when the published packet event is triggered.
func (s *StreamSource) handler(i ...interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.ctx.Err() != nil {
		return
	}
	if len(i) == 2 {
		p, ok := i[1].(*packet.Message)
		if !ok {
			logger.Logger.Error("ListenTopicPublishEvent: type error")
		}

		if p == nil || p.PublishPacket == nil {
			logger.Logger.Error("ListenTopicPublishEvent: publish packet is nil")
			return
		}

		// copy the published packet and set the QoS to QoS0
		var publishPacket = pool.PublishPool.Get()

		pool.CopyPublish(publishPacket, p.PublishPacket)

		publishPacket.QoS = broker.QoS0
		p.PublishPacket = publishPacket

		select {
		case s.messageChan <- p:

		case <-s.ctx.Done():
			s.close()
			return

		default:
			logger.Logger.Warn("drop message")
		}
	}
}
