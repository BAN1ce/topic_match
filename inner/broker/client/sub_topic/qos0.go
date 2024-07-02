package sub_topic

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"go.uber.org/zap"
)

// QoS0 is topic with QoS0
type QoS0 struct {
	ctx           context.Context
	cancel        context.CancelFunc
	topic         string
	messageSource broker.StreamSource
	client        *Client
	meta          *topic.Meta
	messageChan   chan *packet.Message
}

func NewQoS0(meta *topic.Meta, writer client.PacketWriter, messageSource broker.StreamSource) *QoS0 {
	q := &QoS0{
		topic:         meta.Topic,
		meta:          meta,
		messageSource: messageSource,
		client:        NewClient(writer, meta),
	}

	return q
}

// Start starts the QoS0 topic, and it will block until the context is done.
// It will create a publication event listener to listen to the publication event of the store.
func (q *QoS0) Start(ctx context.Context) error {
	q.ctx, q.cancel = context.WithCancel(ctx)
	q.messageChan = make(chan *packet.Message, 10)

	err := q.messageSource.ListenMessage(q.ctx, q.messageChan)
	if err != nil {
		logger.Logger.Error("QoS0Subscriber: listen message error", zap.Error(err))
		return err
	}

	logger.Logger.Debug("QoS0Subscriber: start success", zap.String("topic", q.topic))

	defer func() {
		logger.Logger.Info("QoS0Subscriber: close", zap.String("topic", q.topic),
			zap.String("clientID", q.GetID()))
	}()

	for {
		select {
		case <-q.ctx.Done():
			return nil
		case msg, ok := <-q.messageChan:
			if !ok {
				return nil
			}
			_ = q.publish(msg)

		}
	}
}

// Close closes the QoS0
func (q *QoS0) Close() error {
	if q.ctx == nil {
		return fmt.Errorf("QoS0Subscriber: ctx is nil, not start")
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	if q.cancel != nil {
		q.cancel()
	}
	return nil
}

func (q *QoS0) Publish(publish *packet.Message) error {
	select {
	case q.messageChan <- publish:
		return nil
	default:
		return fmt.Errorf("QoS0Subscriber: publish chan is full")

	}
}

func (q *QoS0) publish(publish *packet.Message) error {
	publish.SetSubIdentifier(byte(q.meta.Identifier))
	publish.PublishPacket.PacketID = q.client.GetPacketWriter().NextPacketID()
	return q.client.Publish(publish)
}

func (q *QoS0) GetID() string {
	return q.client.GetPacketWriter().GetID()
}

func (q *QoS0) GetUnFinishedMessage() []*packet.Message {
	return nil
}

func (q *QoS0) Meta() topic.Meta {
	return *q.meta
}
