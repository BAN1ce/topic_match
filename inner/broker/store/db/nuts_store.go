package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"github.com/google/uuid"
	"github.com/nutsdb/nutsdb"
	"github.com/nutsdb/nutsdb/ds/zset"
	"log"
	"math"
	"strings"
	"time"
)

var (
	nutsDB *nutsdb.DB
)

func initNutsDB(options nutsdb.Options, option ...nutsdb.Option) {
	var (
		err error
	)
	nutsDB, err = nutsdb.Open(
		// nutsdb.DefaultOptions,
		// TODO: support config
		// nutsdb.WithDir("./data/nutsdb"),
		options,
		option...,
	)
	if err != nil {
		log.Fatalln("open db error: ", err)
	}
}

type NutsDB struct {
	db       *nutsdb.DB
	kvBucket string
}

func NewNutsDBStore(options nutsdb.Options, bucketName string, option ...nutsdb.Option) *NutsDB {
	var (
		store = new(NutsDB)
	)
	store.kvBucket = bucketName
	initNutsDB(options, option...)
	store.db = nutsDB
	if store.db == nil {
		logger.Logger.Panic().Msg("local db is nil")
	}
	return store
}

func (s *NutsDB) CreatePacket(topic string, value []byte) (id string, err error) {

	id = uuid.NewString()
	tx, err := s.db.Begin(true)
	if err != nil {
		return "", err
	}
	var timestamp = time.Now().UnixNano()
	if err := tx.ZAdd(topic, []byte(id), float64(timestamp), value); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			logger.Logger.Error().Err(err1).Msg("rollback error")
		}
		return "", err
	} else {
		err = tx.Commit()
	}
	logger.Logger.Debug().Str("topic", topic).Str("value", string(value)).Int64("timestamp", timestamp).Str("messageID", id).Msg("create packet to NutsDB")
	return
}

func (s *NutsDB) ReadFromTimestamp(ctx context.Context, topic string, timestamp time.Time, limit int) ([]*packet.Message, error) {
	var (
		messages []*packet.Message
		err      error
	)
	err = s.db.View(func(tx *nutsdb.Tx) error {
		tmp, err := tx.ZRangeByScore(topic, float64(timestamp.UnixNano()), math.MaxFloat64, &zset.GetByScoreRangeOptions{
			Limit: limit,
		})
		if err != nil {
			return err
		}
		messages = nutsDBValuesBeMessages(tmp, topic)
		return nil
	})
	return messages, err

}

func (s *NutsDB) ReadTopicMessagesByID(ctx context.Context, topic, id string, limit int, include bool) ([]*packet.Message, error) {
	var (
		messages []*packet.Message
		err      error
	)

	err = s.db.View(func(tx *nutsdb.Tx) error {
		score, err := tx.ZScore(topic, []byte(id))
		if err != nil && !errors.Is(err, nutsdb.ErrKeyNotFound) {
			return err
		}

		if score == 0 {
			tmp, err := tx.ZPeekMax(topic)
			if err != nil {
				return err
			} else {
				score = float64(tmp.Score())
			}
		}
		if tmp, err := tx.ZRangeByScore(topic, score, math.MaxFloat64, &zset.GetByScoreRangeOptions{
			Limit: limit,
		}); err != nil {
			return err
		} else {
			if !include {
				for i := 0; i < len(tmp); i++ {
					if string(tmp[i].Key()) == id {
						tmp = append(tmp[:i], tmp[i+1:]...)
						break
					}
				}
			}
			messages = nutsDBValuesBeMessages(tmp, topic)
			return nil
		}
	})
	return messages, err
}

func (s *NutsDB) DeleteTopicMessageID(ctx context.Context, topic, messageID string) error {
	return s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.ZRem(topic, messageID)
	})
}

func (s *NutsDB) DeleteBeforeTime(ctx context.Context, topic string, time time.Time, limit int) error {
	var (
		expiredTotal int
		err          error
		innerError   error
	)

	err = s.db.Update(func(tx *nutsdb.Tx) error {
		expiredTotal, err = tx.ZCount(topic, 0, float64(time.UnixNano()), nil)
		return err
	})

	if err != nil {
		return err
	}

	for i := 0; i < expiredTotal; i += limit {
		var (
			messages []*zset.SortedSetNode
		)
		err = s.db.Update(func(tx *nutsdb.Tx) error {
			messages, innerError = tx.ZRangeByScore(topic, 0, float64(time.UnixNano()), &zset.GetByScoreRangeOptions{
				Limit: limit,
			})
			return errors.Join(err, innerError)
		})

		if err != nil {
			logger.Logger.Error().Err(err).Msg("read before time error")
			continue
		}

		err = s.db.Update(func(tx *nutsdb.Tx) error {
			logger.Logger.Debug().Int("delete message size", len(messages)).Msg("delete before time")
			for _, v := range messages {
				innerError = tx.ZRem(topic, string(v.Key()))
				if err != nil {
					errors.Join(err, innerError)
					logger.Logger.Error().Err(err).Msg("delete before time error")
				}
			}
			return err
		})

	}

	return err
}

func (s *NutsDB) Topics(start, limit int) []string {
	var (
		result = make([]string, 0, limit)
	)
	if err := s.db.View(func(tx *nutsdb.Tx) error {
		if entries, _, err := tx.PrefixScan(s.kvBucket, []byte(""), start, limit); err != nil {
			for _, entry := range entries {
				result = append(result, string(entry.Key))
			}
			return err
		}
		return nil
	}); err != nil {
		if !errors.Is(err, nutsdb.ErrPrefixScan) {
			logger.Logger.Error().Err(err).Msg("range topics error")
		}
	}

	return result
}

func (s *NutsDB) TopicMessageTotal(ctx context.Context, topic string) (int, error) {
	var (
		total int
		err   error
	)
	err = s.db.View(func(tx *nutsdb.Tx) error {
		total, err = tx.ZCard(topic)
		return err
	})
	return total, err

}

func (s *NutsDB) ReadTopicMessage(ctx context.Context, topic string, start, limit int) ([]*packet.Message, error) {
	var (
		messages []*packet.Message
		err      error
	)
	err = s.db.View(func(tx *nutsdb.Tx) error {
		tmp, err := tx.ZRangeByRank(topic, start, limit)
		if err != nil {
			return err
		}
		messages = nutsDBValuesBeMessages(tmp, topic)
		return nil
	})
	return messages, err
}

func nutsDBValuesBeMessages(values []*zset.SortedSetNode, topic string) []*packet.Message {
	var (
		messages []*packet.Message
	)
	for _, v := range values {
		if pubMessage, err := broker.Decode(v.Value); err != nil {
			logger.Logger.Error().Str("topic", topic).Str("messageID", string(v.Key())).Msg("read from NutsDB decode message error")
			continue
		} else {
			messages = append(messages, pubMessage)
		}
	}
	return messages
}

/**

 */
// ---------------------------------------------------------- Key Value MessageStore ----------------------------------------------------------//
/*

 */

func (s *NutsDB) PutKey(ctx context.Context, key, value string) error {
	var err error
	if err = s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Put(s.kvBucket, []byte(key), []byte(value), 0)
	}); err != nil {
		return err
	}
	return nil
}

func (s *NutsDB) ReadKey(ctx context.Context, key string) (string, bool, error) {
	var (
		value string
		err   error
		ok    bool
	)
	err = s.db.View(func(tx *nutsdb.Tx) error {
		if v, err := tx.Get(s.kvBucket, []byte(key)); err != nil {
			// if key not exists, return nil, not error
			// Notice: it's not a good way to check key exists
			// TODO: use other way to check key exists
			if strings.Index(err.Error(), "set not exists") != -1 {
				return nil
			}
			if errors.Is(err, nutsdb.ErrKeyNotFound) {
				err = nil
			}
			return err
		} else {
			value = string(v.Value)
			ok = true
			return nil
		}
	})

	if err != nil {
		return "", false, err
	}
	return value, ok, nil
}

func (s *NutsDB) DeleteKey(ctx context.Context, key string) error {
	if err := s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.Delete(s.kvBucket, []byte(key))
	}); err != nil {
		return err
	}
	return nil
}

func (s *NutsDB) ReadPrefixKey(ctx context.Context, prefix string) (map[string]string, error) {
	var (
		values = make(map[string]string)
		err    error
	)
	defer func() {
	}()
	if err = s.db.View(func(tx *nutsdb.Tx) error {
		// TODO: limit number set 9999, it's dangerous
		if entries, _, err := tx.PrefixScan(s.kvBucket, []byte(prefix), 0, 9999); err != nil {
			return err
		} else {
			for _, entry := range entries {
				values[string(entry.Key)] = string(entry.Value)
			}
			return nil
		}
	}); err != nil {
		if errors.Is(err, nutsdb.ErrPrefixScan) {
			return nil, errors.Join(errs.ErrStoreKeyNotFound, fmt.Errorf("prefix key: %s", prefix))
		}
		return nil, err
	}
	return values, nil
}

func (s *NutsDB) DeletePrefixKey(ctx context.Context, prefix string) error {
	m, err := s.ReadPrefixKey(ctx, prefix)
	if err != nil {
		return err
	}

	for k := range m {
		if tmp := s.DeleteKey(ctx, k); tmp != nil {
			err = errors.Join(err, tmp)
		}
	}
	return err
}

// ---------------------------------------------------------- ZAdd MessageStore ----------------------------------------------------------//

func (s *NutsDB) ZAdd(ctx context.Context, key, member string, score float64) error {
	return s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.ZAdd(s.kvBucket, []byte(key), score, []byte(member))
	})
}

func (s *NutsDB) ZDel(ctx context.Context, key, member string) error {
	return s.db.Update(func(tx *nutsdb.Tx) error {
		return tx.ZRem(s.kvBucket, member)
	})
}

func (s *NutsDB) ZRange(ctx context.Context, key string, start, end float64) ([]string, error) {
	var (
		result []string
	)
	err := s.db.View(func(tx *nutsdb.Tx) error {
		tmp, err := tx.ZRangeByScore(s.kvBucket, start, end, nil)
		if err != nil {
			return err
		}
		for _, v := range tmp {
			result = append(result, v.Key())
		}
		return nil
	})
	return result, err
}

func (s *NutsDB) ZRangeByScore(ctx context.Context, key string, start, end float64) ([]string, error) {
	var (
		result []string
	)
	err := s.db.View(func(tx *nutsdb.Tx) error {
		tmp, err := tx.ZRangeByScore(s.kvBucket, start, end, nil)
		if err != nil {
			return err
		}
		for _, v := range tmp {
			result = append(result, string(v.Key()))
		}
		return nil
	})
	return result, err
}

func (s *NutsDB) Close() error {
	return s.db.Close()
}
