package sub_topic

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"time"
)

type QoS1Option func(q *QoS1)

func QoS1WithLatestMessageID(messageID string) QoS1Option {
	return func(q *QoS1) {
		q.meta.LatestMessageID = messageID
	}
}

type QoS1 struct {
	ctx               context.Context
	cancel            context.CancelFunc
	meta              *topic.Meta
	publishChan       *PublishChan
	client            client.Client
	messageSource     broker.ReadMessageSource
	unfinishedMessage []*packet.Message
	ready             chan struct{}
	qos0Stream        broker.StreamSource
}

func NewQoS1(meta *topic.Meta, writer client.PacketWriter, messageSource broker.ReadMessageSource, qos0 broker.StreamSource, unfinishedMessage []*packet.Message, options ...QoS1Option) *QoS1 {
	q := &QoS1{
		meta:          meta,
		client:        NewQoSWithRetry(NewClient(writer, meta)),
		messageSource: messageSource,
		ready:         make(chan struct{}),
		qos0Stream:    qos0,
	}
	q.unfinishedMessage = unfinishedMessage

	for _, op := range options {
		op(q)
	}
	return q
}

func (q *QoS1) Start(ctx context.Context) error {
	q.ctx, q.cancel = context.WithCancel(ctx)
	if q.meta.WindowSize == 0 {
		// FIXME: config.GetTopic().WindowSize,use client or another config
		q.meta.WindowSize = config.GetTopic().WindowSize
	}

	q.publishChan = NewPublishChan(max(len(q.unfinishedMessage), 1000))
	//q.publishChan = NewPublishChan(0)

	for _, msg := range q.unfinishedMessage {
		q.publishChan.write(msg)
	}

	clear(q.unfinishedMessage)
	q.unfinishedMessage = nil

	close(q.ready)
	if err := q.qos0Stream.ListenMessage(q.ctx, q.publishChan.ch); err != nil {
		err = errors.Join(err, q.qos0Stream.Close())
		return err
	}

	q.listenPublishChan()
	return q.qos0Stream.Close()
}

func (q *QoS1) HandlePublishAck(puback *packets.Puback) (err error) {
	return q.client.HandlePublishAck(puback)
}

func (q *QoS1) listenPublishChan() {
	var (
		delayTime = 5 * time.Second
	)
	defer func() {
		logger.Logger.Debug().Str("topic", q.meta.Topic).Msg("QoS1Subscriber: close listen publish chan")
	}()
	logger.Logger.Debug().Str("topic", q.meta.Topic).Msg("QoS1Subscriber: start listen publish chan")

	for {
		select {
		case <-q.ctx.Done():
			logger.Logger.Info().Str("topic", q.meta.Topic).Msg("QoS1Subscriber: close")
			return

		case msg, ok := <-q.publishChan.ch:
			if !ok {
				return
			}
			// important: the publication topic must be the same as the topic of the subscriber
			logger.Logger.Debug().Str("topic", q.meta.Topic).Msg("QoS1Subscriber: receive from channel")
			if !msg.PublishPacket.Duplicate {
				msg.OwnerClientID = q.client.GetPacketWriter().GetID()
				msg.SetSubIdentifier(byte(q.meta.Identifier))
				if int(msg.PublishPacket.QoS) >= 1 {
					msg.PublishPacket.QoS = broker.QoS1
					msg.PublishPacket.PacketID = q.client.GetPacketWriter().NextPacketID()
				}
				msg.SubscribeTopic = q.meta.Topic
			}

			if err := q.client.Publish(msg); err != nil {
				logger.Logger.Warn().Str("topic", q.meta.Topic).Msg("QoS1Subscriber: publish error")
			}

			if !msg.IsFromSession() && msg.PublishPacket.QoS >= broker.QoS1 {
				q.meta.SetLatestMessageID(msg.MessageID)
			}
		default:
			logger.Logger.Debug().Str("topic", q.meta.Topic).Msg("QoS1Subscriber: start read store")
			message, err := q.messageSource.NextMessages(q.ctx, q.meta.WindowSize, q.meta.LatestMessageID, false, q.publishChan.ch)

			if err != nil && q.ctx.Err() == nil {
				logger.Logger.Error().Err(err).Str("topic", q.meta.Topic).Msg("QoS1Subscriber: read store error")
				time.Sleep(delayTime)
				delayTime *= 2
				if delayTime > 5*time.Minute {
					delayTime = 5 * time.Second
				}
				continue
			}

			delayTime = 5 * time.Second
			for _, m := range message {
				q.publishChan.writeToPublishWithDrop(m)
			}
		}
	}
}

func (q *QoS1) Close() error {
	if q.ctx == nil || q.cancel == nil {
		return fmt.Errorf("QoS1Subscriber: ctx is nil")
	}
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	q.cancel()
	return nil
}

func (q *QoS1) Publish(publish *packet.Message, extra *packet.MessageExtraInfo) error {
	<-q.ready
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	if publish.Duplicate || publish.Retain {
		return q.client.Publish(publish)
	}
	q.publishChan.write(publish)
	return nil
}

func (q *QoS1) GetUnFinishedMessage() []*packet.Message {
	return q.client.GetUnFinishedMessage()
}

func (q *QoS1) Meta() *topic.Meta {
	return q.meta
}
