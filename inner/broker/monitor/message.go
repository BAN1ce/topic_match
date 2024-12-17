package monitor

import (
	"context"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"math"
	"time"
)

type Message struct {
	storeWrapper broker.TopicMessageStore
	tk           time.Ticker
}

func NewMessage(interval time.Duration, store broker.TopicMessageStore) *Message {
	return &Message{
		tk:           *time.NewTicker(interval),
		storeWrapper: store,
	}
}

func (m *Message) Start(ctx context.Context) error {
	for {
		select {
		case <-m.tk.C:
			for i := 0; i <= math.MaxInt; i += 1000 {

				// TODO: add metric to record the delete count and time
				topics := m.storeWrapper.Topics(i, 1000)
				if len(topics) == 0 {
					break
				}

				for _, t := range topics {
					if err := m.storeWrapper.DeleteBeforeTime(ctx, t, time.Now(), 10000); err != nil {
						logger.Logger.Error().Err(err).Str("topic", t).Msg("delete before time error by monitor")
					}
				}
				time.Sleep(3 * time.Second)

			}

		case <-ctx.Done():
			logger.Logger.Info().Msg("store monitor context done")
			return nil
		}
	}

}

func (m *Message) Name() string {
	return "store monitor"
}
