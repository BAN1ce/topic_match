package sub_topic

import (
	"context"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"go.uber.org/zap"
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
}

func (q *QoS2) GetUnfinishedMessage() []*session.UnFinishedMessage {
	//TODO implement me
	panic("implement me")
}

func NewQoS2(meta *topic.Meta, writer client.PacketWriter, messageSource broker.ReadMessageSource, unfinishedMessage []*packet.Message, options ...QoS2Option) *QoS2 {
	t := &QoS2{
		meta:              meta,
		client:            NewQoSWithRetry(NewClient(writer, meta), nil),
		messageSource:     messageSource,
		unfinishedMessage: unfinishedMessage,
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

	q.publishChan = NewPublishChan(max(len(q.unfinishedMessage), q.meta.WindowSize) * 2)
	for _, msg := range q.unfinishedMessage {
		q.publishChan.write(msg)
	}

	clear(q.unfinishedMessage)

	q.listenPublishChan()

	// waiting for exit
	<-ctx.Done()

	if err := q.afterClose(); err != nil {
		logger.Logger.Warn("QoS2Subscriber: after close error = ", zap.Error(err))
	}
	return nil
}

func (q *QoS2) listenPublishChan() {
	var (
		delayTime = 5 * time.Second
	)
	for {
		select {
		case <-q.ctx.Done():
			return

		case msg, ok := <-q.publishChan.ch:
			if !ok {
				return
			}
			msg.SetSubIdentifier(byte(q.meta.Identifier))
			if err := q.client.Publish(msg); err != nil {
				logger.Logger.Warn("QoS2Subscriber: publish error = ", zap.Error(err))
			}
			if !msg.IsFromSession() {
				q.meta.LatestMessageID = msg.MessageID
			}
		default:
			message, err := q.messageSource.NextMessages(q.ctx, q.meta.WindowSize, q.meta.LatestMessageID, false)
			if err != nil {
				logger.Logger.Error("QoS2Subscriber: read store error = ", zap.Error(err), zap.String("store", q.meta.Topic))
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

func (q *QoS2) HandlePublishRec(pubrec *packets.Pubrec) (ok bool, err error) {
	return q.client.HandlePublishRec(pubrec)
}

func (q *QoS2) afterClose() error {
	return nil
}

func (q *QoS2) Publish(publish *packet.Message) error {
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	q.publishChan.write(publish)
	return nil
}

func (q *QoS2) HandlePublishComp(pubcomp *packets.Pubcomp) (bool, error) {
	return q.client.HandelPublishComp(pubcomp)
}

func (q *QoS2) GetUnFinishedMessage() []*packet.Message {
	return q.client.GetUnFinishedMessage()
}

func (q *QoS2) Meta() topic.Meta {
	return *q.meta
}
