package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
)

const (
	SubtreeWriteEvent  = "event.subtree.write"
	SubtreeReadEvent   = "event.subtree.read"
	SubtreeDeleteEvent = "event.subtree.delete"
)

var (
	SubtreeEvent = newSubTree()
)

type Subtree struct {
}

func newSubTree() *Subtree {
	s := &Subtree{}
	return s
}

func (s *Subtree) Emit(event string, data *SubtreeEventData) {
	eventDriver.Emit(event, data)
}

func (s *Subtree) addDefaultListener() {

	eventDriver.AddListener(SubtreeWriteEvent, func(i ...interface{}) {
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

	eventDriver.AddListener(SubtreeReadEvent, func(i ...interface{}) {
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

	eventDriver.AddListener(SubtreeDeleteEvent, func(i ...interface{}) {
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
