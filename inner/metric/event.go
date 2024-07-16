package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CreateListenEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "create_listen_event",
		Help: "create listen event",
	})

	DeleteListenEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "delete_listen_event",
		Help: "delete listen event",
	})

	TriggerListenEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "trigger_listen_event",
		Help: "trigger listen event",
	})
)
