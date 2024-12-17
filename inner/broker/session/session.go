package session

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"github.com/BAN1ce/skyTree/pkg/broker/store"
	"github.com/BAN1ce/skyTree/pkg/broker/topic"
	"github.com/BAN1ce/skyTree/pkg/errs"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"time"
)

type Session struct {
	clientID string
	store    *store.KeyValueStoreWithTimeout
}

func newSession(ctx context.Context, clientID string, keyStore store.KeyStore) *Session {
	s := &Session{
		clientID: clientID,
		store: store.NewKeyValueStoreWithTimout(
			keyStore,
			3*time.Second),
	}
	if err := s.store.PutKey(ctx, clientSessionCreatedTimeKey(s.clientID), time.Now().String()); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("put key failed")
	}
	return s
}

func (s *Session) Release() {
	// delete session prefix key  like session/client/xxx
	if err := s.store.DefaultDeleteKey(clientSessionPrefix(s.clientID)); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("release session failed")
	}
}

// ----------------------------------------------------------------- Sub topic ----------------------------------------------------------------- //

// ReadSubTopics returns all sub topics of the client.
func (s *Session) ReadSubTopics() (topics []*topic.Meta) {
	m, err := s.store.DefaultReadPrefixKey(context.TODO(), clientSubTopicKeyPrefix(s.clientID))
	if err != nil && !errors.Is(err, errs.ErrStoreKeyNotFound) {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("read sub topics failed")
		return nil
	}
	if len(m) == 0 {
		logger.Logger.Debug().Str("clientID", s.clientID).Msg("no sub topics")
	}
	for _, v := range m {
		var meta *topic.Meta
		// QoS should not greater than 2, so int is enough
		if err := json.Unmarshal([]byte(v), &meta); err != nil {
			logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("unmarshal sub option failed")
			continue
		}

		topics = append(topics, meta)
	}
	return
}

// CreateSubTopic creates a sub topic for the client. store the sub option to the session.
func (s *Session) CreateSubTopic(meta *topic.Meta) {
	var (
		topicName = meta.Topic
	)
	if topicName == "" {
		logger.Logger.Error().Str("clientID", s.clientID).Msg("topicName should not be empty")
		return
	}
	data, _ := json.Marshal(meta)
	if err := s.store.DefaultPutKey(clientSubTopicKey(s.clientID, topicName), string(data)); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Str("topicName", topicName).Int32("qos", meta.QoS).Msg("create sub topic failed")
	}
}

// DeleteSubTopic deletes a sub topic for the client.
func (s *Session) DeleteSubTopic(topic string) {
	if topic == "" {
		return
	}
	if err := s.store.DefaultDeleteKey(clientSubTopicKey(s.clientID, topic)); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Str("topic", topic).Msg("delete sub topic failed")
	}
}

// ----------------------------------------------------------------- topic Latest Pushed Message ----------------------------------------------------------------- //

func (s *Session) ReadTopicLatestPushedMessageID(topic string) (messageID string, ok bool) {
	id, _, err := s.store.DefaultReadKey(clientLatestAckedMessageKey(s.clientID, topic))
	if err != nil {
		logger.Logger.Error().Str("clientID", s.clientID).Err(err).Msg("read topic last acked message id failed")
		return "", false
	}
	return id, true
}

func (s *Session) SetTopicLatestPushedMessageID(topic string, messageID string) {
	if err := s.store.DefaultPutKey(clientLatestAckedMessageKey(s.clientID, topic), messageID); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Str("topic", topic).Str("messageID", messageID).Msg("set topic last acked message id failed")
	}
}

func (s *Session) DeleteTopicLatestPushedMessageID(topic string, messageID string) {
	if err := s.store.DefaultDeleteKey(clientLatestAckedMessageKey(s.clientID, topic)); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Str("topic", topic).Str("messageID", messageID).Msg("delete topic last acked message id failed")
	}
}

// ----------------------------------------------------------------- topic UnFinished Message ----------------------------------------------------------------- //

func (s *Session) CreateTopicUnFinishedMessage(topic string, message []*packet.Message) {
	prefix := clientUnfinishedMessageKey(s.clientID, topic)
	for _, m := range message {
		payload, err := json.Marshal(m)
		if err != nil {
			logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("marshal unfinished message failed")
			continue
		}
		if err := s.store.DefaultPutKey(prefix+"/"+m.MessageID, string(payload)); err != nil {
			logger.Logger.Error().Err(err).Str("clientID", s.clientID).Str("topic", topic).Str("messageID", m.MessageID).Msg("create topic unfinished message failed")
		}
	}
}

func (s *Session) ReadTopicUnFinishedMessage(topic string) (message []*packet.Message) {
	prefix := clientUnfinishedMessageKey(s.clientID, topic)

	m, err := s.store.DefaultReadPrefixKey(context.TODO(), prefix)
	if errors.Is(err, errs.ErrStoreKeyNotFound) {
		return nil
	}

	if err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("read topic unfinished message failed")
		return
	}

	for _, v := range m {
		var m *packet.Message
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("unmarshal unfinished message failed")
			continue
		}
		if m != nil {
			message = append(message, m)
		} else {
			logger.Logger.Warn().Str("clientID", s.clientID).Str("topic", topic).Str("body", string(v)).Msg("json unmarshal message is nil")
		}
	}
	return message
}

func (s *Session) DeleteTopicUnFinishedMessage(topic string, _ string) {
	if err := s.store.DefaultDeleteKey(clientUnfinishedMessageKey(s.clientID, topic)); err != nil {
		logger.Logger.Error().Str("clientID", s.clientID).Str("topic", topic).Msg("delete topic unfinished message failed")
	}
}

// ----------------------------------------------------------------- Will Message ----------------------------------------------------------------- //

func (s *Session) GetWillMessage() (*session.WillMessage, bool, error) {
	var message session.WillMessage
	value, ok, err := s.store.DefaultReadKey(clientWillMessageKey(s.clientID))
	if ok {
		if err = json.Unmarshal([]byte(value), &message); err != nil {
			logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("unmarshal will message failed")
		}
	}

	return &message, ok, err
}

func (s *Session) SetWillMessage(message *session.WillMessage) error {
	message.CreatedTime = time.Now().String()
	if jBody, err := json.Marshal(message); err != nil {
		return err
	} else {
		logger.Logger.Debug().Str("clientID", s.clientID).Msg("set will message")
		return s.store.DefaultPutKey(clientWillMessageKey(s.clientID), string(jBody))
	}
}

func (s *Session) DeleteWillMessage() error {
	return s.store.DefaultDeleteKey(clientWillMessageKey(s.clientID))
}

// ----------------------------------------------------------------- Connect Properties ----------------------------------------------------------------- //

func (s *Session) GetConnectProperties() (*session.ConnectProperties, error) {
	var (
		properties session.ConnectProperties
	)
	value, ok, err := s.store.DefaultReadKey(clientConnectPropertiesKey(s.clientID))
	if err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("read connect properties failed")
		return &properties, err
	}
	if !ok {
		return &properties, errs.ErrSessionConnectPropertiesNotFound
	}
	if err := json.Unmarshal([]byte(value), &properties); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("unmarshal connect properties failed")
		return &properties, err
	}
	return &properties, nil
}

func (s *Session) SetConnectProperties(properties *session.ConnectProperties) error {
	value, err := json.Marshal(properties)
	if err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("marshal connect properties failed")
		return err
	}
	if err := s.store.DefaultPutKey(clientConnectPropertiesKey(s.clientID), string(value)); err != nil {
		logger.Logger.Error().Err(err).Str("clientID", s.clientID).Msg("set connect properties failed")
		return err
	}
	return nil
}

func (s *Session) SetExpiryInterval(i int64) {
	//TODO implement me
	panic("implement me")
}

func (s *Session) GetExpiryInterval() int64 {
	// TODO: implement me
	return 0
}
