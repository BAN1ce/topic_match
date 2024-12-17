package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	DropQoS0MessageCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "drop_qos0_message_count",
	})
)
