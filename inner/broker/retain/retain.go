package retain

import (
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker/retain"
	"github.com/BAN1ce/skyTree/pkg/broker/store"
	"time"
)

type DB struct {
	store *store.KeyValueStoreWithTimeout
}

func NewRetainDB(keyStore store.KeyStore) *DB {
	return &DB{store: store.NewKeyValueStoreWithTimout(keyStore, 3*time.Second)}
}

func (d *DB) PutRetainMessage(message *retain.Message) error {
	logger.Logger.Debug().Str("topic", message.Topic).Str("body", string(message.Encode())).Msg("put retain message")
	//if len(message.Payload) == 0 {
	//	return d.store.DefaultDeleteKey(store.GetTopicRetainKey(message.Topic))
	//}

	return d.store.DefaultPutKey(store.GetTopicRetainKey(message.Topic), string(message.Encode()))
}

func (d *DB) GetRetainMessage(topic string) (*retain.Message, bool) {
	data, ok, err := d.store.DefaultReadKey(store.GetTopicRetainKey(topic))
	if err != nil {
		logger.Logger.Info().Err(err).Str("topic", topic).Msg("get retain message failed")
		return nil, false
	}
	if !ok {
		return nil, false
	}
	message := &retain.Message{}
	err = message.Decode([]byte(data))
	if err != nil {
		return nil, true
	}
	return message, true
}

func (d *DB) DeleteRetainMessage(topic string) error {
	return d.store.DefaultDeleteKey(store.GetTopicRetainKey(topic))
}
