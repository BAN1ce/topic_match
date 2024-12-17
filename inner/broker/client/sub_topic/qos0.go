package sub_topic

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/BAN1ce/skyTree/pkg/pool"
)

// QoS0 is topic with QoS0
type QoS0 struct {
	ctx           context.Context
	cancel        context.CancelFunc
	messageSource broker.StreamSource
	client        *Client
	meta          *topic.Meta
	messageChan   chan *packet.Message
	ready         chan struct{}
}

func NewQoS0(meta *topic.Meta, writer client.PacketWriter, messageSource broker.StreamSource) *QoS0 {
	q := &QoS0{
		meta:          meta,
		messageSource: messageSource,
		client:        NewClient(writer, meta),
		ready:         make(chan struct{}),
	}
	q.messageChan = make(chan *packet.Message, 1000)

	return q
}

// Start starts the QoS0 topic, and it will block until the context is done.
// It will create a publication event listener to listen to the publication event of the store.
func (q *QoS0) Start(ctx context.Context) error {
	q.ctx, q.cancel = context.WithCancel(ctx)

	err := q.messageSource.ListenMessage(q.ctx, q.messageChan)
	if err != nil {
		logger.Logger.Error().Err(err).Msg("QoS0Subscriber: listen message error")
		return err
	}

	defer func() {
		logger.Logger.Info().Str("topic", q.meta.Topic).Str("clientID", q.GetID()).Msg("QoS0Subscriber: close")
	}()

	close(q.ready)
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
	logger.Logger.Info().Str("topic", q.meta.Topic).Msg("QoS0Subscriber: close")
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

func (q *QoS0) Publish(publish *packet.Message, extra *packet.MessageExtraInfo) error {
	<-q.ready

	select {
	case q.messageChan <- publish:
		return nil
	case <-q.ctx.Done():
		return fmt.Errorf("QoS0Subscriber: ctx is done")
	default:
		return fmt.Errorf("QoS0Subscriber: publish chan is full")

	}
}

func (q *QoS0) publish(message *packet.Message) error {
	var (
		publish    = pool.PublishPool.Get()
		newMessage = &packet.Message{}
	)
	defer pool.PublishPool.Put(publish)

	pool.CopyPublish(publish, message.PublishPacket)

	newMessage.SubscribeTopic = q.meta.Topic
	newMessage.PublishPacket = publish
	newMessage.PublishPacket.QoS = broker.QoS0
	newMessage.SetSubIdentifier(byte(q.meta.Identifier))
	newMessage.PublishPacket.PacketID = q.client.GetPacketWriter().NextPacketID()

	return q.client.Publish(newMessage)
}

func (q *QoS0) GetID() string {
	return q.client.GetPacketWriter().GetID()
}

func (q *QoS0) GetUnFinishedMessage() []*packet.Message {
	return nil
}

func (q *QoS0) Meta() *topic.Meta {
	return q.meta
}
