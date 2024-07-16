package sub_topic

import (
	"context"
	"fmt"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/client"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/eclipse/paho.golang/packets"
	"go.uber.org/zap"
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
	meta              topic.Meta
	publishChan       *PublishChan
	client            client.Client
	messageSource     broker.ReadMessageSource
	unfinishedMessage []*packet.Message
}

func NewQoS1(meta *topic.Meta, writer client.PacketWriter, messageSource broker.ReadMessageSource, unfinishedMessage []*packet.Message, options ...QoS1Option) *QoS1 {
	q := &QoS1{
		meta:          *meta,
		client:        NewQoSWithRetry(NewClient(writer, meta), nil),
		messageSource: messageSource,
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
	q.publishChan = NewPublishChan(max(len(q.unfinishedMessage), q.meta.WindowSize) * 2)
	//q.ch = make(chan *packet.Message, 0)

	for _, msg := range q.unfinishedMessage {
		q.publishChan.write(msg)
	}
	clear(q.unfinishedMessage)
	q.unfinishedMessage = nil
	q.listenPublishChan()
	return nil
}

func (q *QoS1) HandlePublishAck(puback *packets.Puback) (ok bool, err error) {
	return q.client.HandlePublishAck(puback)
}

func (q *QoS1) listenPublishChan() {
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
			msg.PublishPacket.QoS = broker.QoS1
			msg.PublishPacket.PacketID = q.client.GetPacketWriter().NextPacketID()
			msg.FromTopic = q.meta.Topic
			if err := q.client.Publish(msg); err != nil {
				logger.Logger.Warn("QoS1Subscriber: publish error = ", zap.Error(err))
			}
			if !msg.IsFromSession() {
				q.meta.LatestMessageID = msg.MessageID
			}
		default:
			logger.Logger.Debug("QoS1Subscriber: start read store")
			message, err := q.messageSource.NextMessages(q.ctx, q.meta.WindowSize, q.meta.LatestMessageID, false)
			if err != nil {
				logger.Logger.Error("QoS1Subscriber: read store error = ", zap.Error(err), zap.String("store", q.meta.Topic))
				time.Sleep(delayTime)
				delayTime *= 2
				if delayTime > 5*time.Minute {
					delayTime = 5 * time.Second
				}
				continue
			}
			logger.Logger.Debug("QoS1Subscriber: read store message", zap.Any("message", message))

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

func (q *QoS1) Publish(publish *packet.Message) error {
	if q.ctx.Err() != nil {
		return q.ctx.Err()
	}
	q.publishChan.write(publish)

	return nil
}

func (q *QoS1) GetUnFinishedMessage() []*packet.Message {
	return q.client.GetUnFinishedMessage()
}

func (q *QoS1) Meta() topic.Meta {
	return q.meta
}
