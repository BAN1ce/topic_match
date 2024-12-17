package share

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/packet"
	"sync"
)

// MessageSource is the source of the message
// We can get the messages of one topic from the message source
type MessageSource struct {
	cancel  context.CancelFunc
	readMux sync.RWMutex

	topic          string
	fullShareTopic string

	queue           *Queue // Element: *packet.Message
	latestMessageID string

	// messageSource is the source of the message
	// it's a store that stored all the messages
	messageSource broker.ReadMessageSource
}

func newMessageSource(fullShareTopic, topic string, source broker.ReadMessageSource) *MessageSource {
	return &MessageSource{
		fullShareTopic: fullShareTopic,
		topic:          topic,
		queue:          NewQueue(),
		messageSource:  source,
	}
}

func (m *MessageSource) NextMessages(ctx context.Context, n int, startMessageID string, include bool, in chan *packet.Message) ([]*packet.Message, error) {
	return m.nextMessage(ctx, n)
}

func (m *MessageSource) ListenMessage(ctx context.Context, writer chan *packet.Message) error {
	var (
		size = 10
	)
	ctx, m.cancel = context.WithCancel(ctx)

	go func(ctx context.Context, c chan *packet.Message, size int) {
		defer close(c)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if message, err := m.nextMessage(ctx, size); err != nil {
					return
				} else {
					for i := 0; i < len(message); i++ {
						select {
						case c <- message[i]:

						case <-ctx.Done():
							return

						}
					}
				}
			}
		}
	}(ctx, writer, size)

	return nil
}

func (m *MessageSource) Close() error {
	if m.cancel != nil {
		m.cancel()
	}
	logger.Logger.Info().Str("topic", m.topic).Str("shareTopic", m.fullShareTopic).Msg("share topic message source close")
	return nil
}

func (m *MessageSource) nextMessage(ctx context.Context, n int) (message []*packet.Message, err error) {
	// fill the share topic for every message
	defer m.fillShareTopic(message)

	m.readMux.RLock()
	if m.queue.Len() != 0 {
		defer m.readMux.RUnlock()
		message = m.queue.PopQueue(n)
		return
	}
	m.readMux.RUnlock()

	m.readMux.Lock()
	defer m.readMux.Unlock()

	if m.queue.Len() != 0 {
		message = m.queue.PopQueue(n)
		return
	}

	message, err = m.messageSource.NextMessages(ctx, 10, m.latestMessageID, false, nil)
	if err != nil {
		return nil, err
	}

	if len(message) == 0 {
		return nil, nil
	}

	m.queue.AppendMessage(message)
	m.latestMessageID = message[len(message)-1].MessageID
	message = m.queue.PopQueue(n)
	return
}

func (m *MessageSource) fillShareTopic(message []*packet.Message) []*packet.Message {
	for i := 0; i < len(message); i++ {
		message[i].ShareTopic = m.fullShareTopic
	}
	return message
}
