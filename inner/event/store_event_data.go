package event

import "time"

type StoreEventData struct {
	Topic     string
	MessageID string
	Success   bool
	Duration  time.Duration
	Count     int
	QoS       int
}

func (s *StoreEventData) name() {

}
