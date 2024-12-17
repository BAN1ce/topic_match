package sub_topic

import (
	"context"
	"errors"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"time"
)

type QoS2Option func(q *QoS2)

func QoS2WithLatestMessageID(messageID string) QoS2Option {
	return func(q *QoS2) {
		q.meta.LatestMessageID = messageID
	}
}

type QoS2 struct {
	ctx               context.Context
	cancel            context.CancelFunc
	meta              *topic.Meta
	client            client.Client
	publishChan       *PublishChan
	messageSource     broker.ReadMessageSource
	subOption         *packets.SubOptions
	unfinishedMessage []*packet.Message
	ready             chan struct{}
	qos0Stream        broker.StreamSource
}

func (q *QoS2) GetUnfinishedMessage() []*session.UnFinishedMessage {
	//TODO implement me
	panic("implement me")
}

func NewQoS2(meta *topic.Meta, writer client.PacketWriter, messageSource broker.ReadMessageSource, qos0 broker.StreamSource, unfinishedMessage []*packet.Message, options ...QoS2Option) *QoS2 {
	t := &QoS2{
		meta:              meta,
		client:            NewQoSWithRetry(NewClient(writer, meta)),
		messageSource:     messageSource,
		unfinishedMessage: unfinishedMessage,
		ready:             make(chan struct{}),
		qos0Stream:        qos0,
	}
	for _, op := range options {
		op(t)
	}
	return t
}

func (q *QoS2) Start(ctx context.Context) error {
	q.ctx, q.cancel = context.WithCancel(ctx)
	if q.meta.WindowSize == 0 {
		// FIXME: config.GetTopic().WindowSize,use client or another config
		q.meta.WindowSize = config.GetTopic().WindowSize
	}

	q.publishChan = NewPublishChan(max(len(q.unfinishedMessage), 2000))
	for _, msg := range q.unfinishedMessage {
		q.publishChan.write(msg)
	}

	clear(q.unfinishedMessage)

	close(q.ready)

	if err := q.qos0Stream.ListenMessage(q.ctx, q.publishChan.ch); err != nil {
		err = errors.Join(err, q.qos0Stream.Close())
		return err
	}

	q.listenPublishChan()

	// waiting for exit
	<-ctx.Done()

	if err := q.afterClose(); err != nil {
		logger.Logger.Warn().Err(err).Msg("QoS2Subscriber: after close error")
	}
	return q.qos0Stream.Close()
}

func (q *QoS2) listenPublishChan() {
	var (
		delayTime = 5 * time.Second
	)
	for {
		select {
		case <-q.ctx.Done():
			logger.Logger.Info().Str("topic", q.meta.Topic).Msg("QoS2Subscriber: close")
			return

		case msg, ok := <-q.publishChan.ch:
			if !ok {
				return
			}

			if msg.PublishPacket.QoS != 0 {
				msg.PublishPacket.PacketID = q.client.GetPacketWriter().NextPacketID()
			}

			if !msg.PublishPacket.Duplicate {
				msg.OwnerClientID = q.client.GetPacketWriter().GetID()
				msg.SetSubIdentifier(byte(q.meta.Identifier))
				msg.SubscribeTopic = q.meta.Topic
			}

			if err := q.client.Publish(msg); err != nil {
				logger.Logger.Warn().Err(err).Msg("QoS2Subscriber: publish error")
			}
			if !msg.IsFromSession() {
				q.meta.SetLatestMessageID(msg.MessageID)
			}

		default:
			message, err := q.messageSource.NextMessages(q.ctx, q.meta.WindowSize, q.meta.LatestMessageID, false, q.publishChan.ch)
			if err != nil && q.ctx.Err() == nil {
				logger.Logger.Error().Str("topic", q.meta.Topic).Err(err).Msg("QoS2Subscriber: read store error")
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

func (q *QoS2) Close() error {
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	q.cancel()
	return nil
}

func (q *QoS2) HandlePublishRec(pubrec *packets.Pubrec) (err error) {
	return q.client.HandlePublishRec(pubrec)
}

func (q *QoS2) afterClose() error {
	return nil
}

func (q *QoS2) Publish(publish *packet.Message, extra *packet.MessageExtraInfo) error {
	<-q.ready

	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	if publish.PublishPacket.Duplicate || publish.Retain {
		return q.client.Publish(publish)
	}

	q.publishChan.write(publish)
	return nil
}

func (q *QoS2) HandlePublishComp(pubcomp *packets.Pubcomp) error {
	return q.client.HandelPublishComp(pubcomp)
}

func (q *QoS2) GetUnFinishedMessage() []*packet.Message {
	return q.client.GetUnFinishedMessage()
}

func (q *QoS2) Meta() *topic.Meta {
	return q.meta.Meta()
}
