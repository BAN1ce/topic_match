package event

import "time"

type SubtreeEventData struct {
	Duration     time.Duration
	Success      bool
	NoSubscriber bool
}
