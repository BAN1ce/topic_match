package message

import (
	"bytes"
	"context"
	"errors"
	"github.com/BAN1ce/skyTree/config"
	"github.com/BAN1ce/skyTree/inner/broker/store"
	event2 "github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	packet2 "github.com/BAN1ce/skyTree/pkg/packet"
	"time"
)

// Wrapper is a wrapper for the store
// it helps to store the message after encode to bytes to the store
// and set the expired time for the message
// it will emit the store event
// and handle the store error
type Wrapper struct {
	store broker.TopicMessageStore
	event broker.MessageStoreEvent
}

func (w *Wrapper) DeleteBeforeTime(ctx context.Context, topic string, time time.Time, limit int) error {
	err := w.store.DeleteBeforeTime(ctx, topic, time, limit)
	return err
}

func NewStoreWrapper(store broker.TopicMessageStore, event broker.MessageStoreEvent) *Wrapper {
	return &Wrapper{
		store: store,
		event: event,
	}
}

// StorePublishPacket stores the published packet to store
// topics include the origin topic and the topic of the wildcard subscription
// if qos is 0, the message will not store
func (w *Wrapper) StorePublishPacket(topics map[string]int32, packet *packet2.Message) (messageID string, err error) {
	if len(topics) == 0 {
		logger.Logger.Warn().Msg("store publish packet with empty topics, should not be empty")
		return "", nil
	}

	var (
		// there doesn't use bytes.BufferPool, because the store maybe async
		encodedData = bytes.NewBuffer(nil)
		expiredTime int64
		errs        error
	)

	// if the message has the message expiry property, set the expired time
	// if not set the default expired time
	if packet.PublishPacket.Properties != nil && packet.PublishPacket.Properties.MessageExpiry != nil {
		expiredTime = time.Now().Add(time.Duration(int64(*packet.PublishPacket.Properties.MessageExpiry)) * time.Second).Unix()
	} else {
		// default expired time
		expiredTime = time.Now().Add(time.Duration(config.GetConfig().Store.MessageExpired) * time.Second).Unix()
	}

	// set the expired time for the message
	packet.ExpiredTime = expiredTime

	// publish packet encode to bytes
	if err = broker.Encode(store.DefaultSerializerVersion, packet, encodedData); err != nil {
		return "", err
	}

	for topic, qos := range topics {
		if qos == 0 {
			// qos 0 message will not store
			continue
		}

		// store message bytes
		messageID, err = w.CreatePacket(topic, encodedData.Bytes())

		if err != nil {
			errs = errors.Join(errs, err)
			logger.Logger.Error().Err(err).Str("topic", topic).Msg("create packet to store error")
		} else {
			logger.Logger.Debug().Str("topic", topic).Str("message_id", messageID).Str("payload", string(packet.PublishPacket.Payload)).Msg("create packet to store success")
		}

	}
	err = errs

	return messageID, err
}

func (w *Wrapper) DeleteTopicMessageID(ctx context.Context, topic, messageID string) error {
	start := time.Now()
	err := w.store.DeleteTopicMessageID(ctx, topic, messageID)
	w.event.EmitDelete(&event2.StoreEventData{
		Success:  err == nil,
		Topic:    topic,
		Duration: time.Since(start),
	})
	return err
}

func (w *Wrapper) ReadFromTimestamp(ctx context.Context, topic string, timestamp time.Time, limit int) ([]*packet2.Message, error) {
	start := time.Now()
	result, err := w.store.ReadFromTimestamp(ctx, topic, timestamp, limit)
	w.event.EmitRead(&event2.StoreEventData{
		Topic:    topic,
		Success:  err == nil,
		Duration: time.Since(start),
		Count:    len(result),
	})
	return result, err
}

// ReadTopicMessagesByID read the messages belong the topic and the message id from the message store
// if the message is expired, delete the message, and return the messages after filter
func (w *Wrapper) ReadTopicMessagesByID(ctx context.Context, topic, id string, limit int, include bool) ([]*packet2.Message, error) {
	var (
		start          = time.Now()
		filteredResult []*packet2.Message
	)
	result, err := w.store.ReadTopicMessagesByID(ctx, topic, id, limit, include)
	w.event.EmitRead(&event2.StoreEventData{
		Topic:    topic,
		Success:  err == nil,
		Duration: time.Since(start),
		Count:    len(result),
	})

	for _, m := range result {
		if m.IsExpired() {
			if err1 := w.store.DeleteTopicMessageID(ctx, topic, m.MessageID); err1 != nil {
				logger.Logger.Error().Err(err1).Str("topic", topic).Str("message_id", m.MessageID).Msg("delete expired message error")
			}
		} else {
			filteredResult = append(filteredResult, m)
		}
	}

	return filteredResult, err
}

func (w *Wrapper) CreatePacket(topic string, value []byte) (id string, err error) {
	start := time.Now()
	id, err = w.store.CreatePacket(topic, value)
	w.event.EmitStored(&event2.StoreEventData{
		Topic:     topic,
		MessageID: id,
		Success:   err == nil,
		Duration:  time.Since(start),
	})
	return
}

// Topics return the topics from the message store
func (w *Wrapper) Topics(start, limit int) []string {
	return w.store.Topics(start, limit)
}

// TopicMessageTotal return the total message count of the topic in the message store
func (w *Wrapper) TopicMessageTotal(ctx context.Context, topic string) (int, error) {
	return w.store.TopicMessageTotal(ctx, topic)
}

// ReadTopicMessage read the messages belong the topic from the message store
func (w *Wrapper) ReadTopicMessage(ctx context.Context, topic string, start, limit int) ([]*packet2.Message, error) {
	return w.store.ReadTopicMessage(ctx, topic, start, limit)
}
