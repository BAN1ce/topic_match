package state

import (
	"context"
	"errors"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/store"
	"github.com/nutsdb/nutsdb"
	"time"
)

type State struct {
	store *store.KeyValueStoreWithTimeout
}

func NewState(keyStore store.KeyStore) *State {
	return &State{
		store: store.NewKeyValueStoreWithTimout(keyStore, 3*time.Second),
	}
}

func (s *State) ReadRetainMessageID(topic string) ([]string, error) {
	var (
		messageIDs []string
		value, err = s.store.DefaultReadPrefixKey(context.TODO(), broker.TopicRetainMessage(topic).String())
	)
	if err != nil {
		return messageIDs, err
	}
	for _, v := range value {
		messageIDs = append(messageIDs, v)
	}
	return messageIDs, err
}

func (s *State) CreateRetainMessageID(topic, messageID string) error {
	return s.store.DefaultPutKey(broker.TopicRetainMessage(topic).String(), messageID)
}

func (s *State) DeleteRetainMessageID(topic string) error {
	return s.store.DefaultDeleteKey(broker.TopicRetainMessage(topic).String())
}

const (
	WillMessageZSet = "will_message_zset"
)

// ----------------- will message  z set -----------------//

func (s *State) AddDelayWillTask(ctx context.Context, clientID string, triggerTime time.Time) error {
	var (
		errs error
	)

	err := s.store.ZAdd(ctx, clientID, "", float64(triggerTime.Unix()))
	if err != nil {
		errs = errors.Join(errs, err)
		logger.Logger.Error().Err(err).Msg("add will message for adding will message error")
		err = s.store.DeleteKey(ctx, clientID)
		if err != nil {
			errs = errors.Join(errs, err)
			logger.Logger.Error().Err(err).Msg("delete will message for adding will message error")
		}
	}

	return errs
}

func (s *State) GetWillClientID(ctx context.Context, time time.Time) []string {
	clientID, err := s.store.ZRangeByScore(ctx, WillMessageZSet, 0, float64(time.Unix()))
	if err != nil && !errors.Is(err, nutsdb.ErrBucket) {
		logger.Logger.Error().Err(err).Msg("get will message error")
		return nil
	}

	return clientID

}

func (s *State) DeleteWillClientID(ctx context.Context, key []string) {
	for _, v := range key {
		if err := s.store.ZDel(ctx, WillMessageZSet, v); err != nil {
			logger.Logger.Error().Err(err).Str("key", v).Msg("delete will message error")
		}
	}

}

func (s *State) Close() error {
	return s.store.Close()
}
