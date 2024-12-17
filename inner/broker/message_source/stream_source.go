package message_source

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"sync"
)

type StreamSource struct {
	ctx             context.Context
	cancel          context.CancelFunc
	topic           string
	publishListener broker.HandlePublishDoneEvent
	messageChan     chan *packet.Message
	mux             sync.RWMutex
	maxQoS          int
}

// NewStreamSource creates a new event source with the given topic.
func NewStreamSource(topic string, listener broker.HandlePublishDoneEvent, maxQoS int) *StreamSource {
	return &StreamSource{
		topic:           topic,
		publishListener: listener,
		maxQoS:          maxQoS,
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
	s.publishListener.DeleteListenPublishEvent(s.topic, s.handler)
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

		//subTopic,ok := i[0].(string)

		p, ok := i[1].(*packet.Message)
		if !ok {
			logger.Logger.Error().Msg("ListenTopicPublishEvent: type error")
		}

		if p == nil || p.PublishPacket == nil {
			logger.Logger.Error().Msg("ListenTopicPublishEvent: publish packet is nil")
			return
		}

		if int(p.PublishPacket.QoS) > s.maxQoS {
			// ignore the message
			// the message is not the qos level that the subscriber wants to listen to
			return
		}

		select {
		case s.messageChan <- p:

		case <-s.ctx.Done():
			s.close()
			return

		default:
			metric.DropQoS0MessageCount.Inc()
		}
	}
}
