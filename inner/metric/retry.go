package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	PublishRetryTaskCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "mqtt_publish_retry_task_count",
		Help: "mqtt publish retry task count",
	}, []string{"topic"})

	PublishRetryTaskCurrent = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mqtt_publish_retry_task_current",
	})

	PublishRetryTimeout = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_publish_retry_timeout",
	})

	PublishRetryDelete = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_publish_retry_deleted",
	})

	ExecuteRetryTask = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_execute_retry_task",
	})
)

//var (
//
//	ReceiveQoS2Publish = promauto.NewCounter(prometheus.CounterOpts{
//
//	}
//)
