package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/kataras/go-events"
	"time"
)

const (
	SubtreeWriteEvent  = "event.subtree.write"
	SubtreeReadEvent   = "event.subtree.read"
	SubtreeDeleteEvent = "event.subtree.delete"
)

var (
	SubtreeEvent = newSubTree()
)

type SubtreeEventData struct {
	Duration     time.Duration
	Success      bool
	NoSubscriber bool
}

type Subtree struct {
}

func newSubTree() *Subtree {
	s := &Subtree{}
	s.registerMetric()
	return s
}

func (s *Subtree) Emit(event string, data *SubtreeEventData) {
	Driver.Emit(events.EventName(event), data)
}

func (s *Subtree) registerMetric() {

	Driver.AddListener(SubtreeWriteEvent, func(i ...interface{}) {
		if len(i) == 0 {
			return
		}
		data, ok := i[0].(*SubtreeEventData)
		if !ok {
			return
		}
		metric.SubTreeWriteHistogram.Observe(data.Duration.Seconds())
		if !data.Success {
			metric.SubTreeWriteFailed.Inc()
		}
	})

	Driver.AddListener(SubtreeReadEvent, func(i ...interface{}) {
		if len(i) == 0 {
			return
		}
		data, ok := i[0].(*SubtreeEventData)
		if !ok {
			return
		}
		metric.SubTreeReadHistogram.Observe(data.Duration.Seconds())
		if !data.Success {
			metric.SubTreeReadFailed.Inc()
		}

		metric.SubTreeReadNoSubscriber.Inc()
	})

	Driver.AddListener(SubtreeDeleteEvent, func(i ...interface{}) {
		if len(i) == 0 {
			return
		}
		data, ok := i[0].(*SubtreeEventData)
		if !ok {
			return
		}
		metric.SubTreeDeleteHistogram.Observe(data.Duration.Seconds())
		if !data.Success {
			metric.SubTreeDeleteFailed.Inc()
		}
	})

}
