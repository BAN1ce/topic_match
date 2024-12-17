package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// MessageStore Read

var (
	StoreReadMessageCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_read_message_count",
		Help: "The total number of messages read from the store",
	})

	StoreReadRequestFailedCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_read_request_failed_count",
		Help: "The total number of failed read requests to the store",
	})

	StoreReadDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "store_read_duration",
		Help: "The duration of read requests to the store",
		Buckets: []float64{
			0.05, 0.1, 0.5, 1, 2, 5, 10,
		},
	})

	StoreReadRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_read_request_count",
		Help: "The total number of read requests to the store",
	})
)

// MessageStore Save

var (
	StoreWriteRequestFailedCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_write_request_failed_count",
		Help: "The total number of failed store requests to the store",
	})

	StoreWriteRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_write_request_count",
		Help: "The total number of store requests to the store",
	})

	StoreWriteDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "store_write_duration",
		Help: "The duration of store requests to the store",
		Buckets: []float64{
			0.05, 0.1, 0.5, 1, 2, 5, 10,
		},
	})
)

// MessageStore Delete metric

var (
	StoreDeleteRequestFailedCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_delete_request_failed_count",
		Help: "The total number of failed delete requests to the store",
	})

	StoreDeleteDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "store_delete_duration",
		Help: "The duration of delete requests to the store",
		Buckets: []float64{
			0.05, 0.1, 0.5, 1, 2, 5, 10,
		},
	})

	StoreDeleteRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_delete_request_count",
		Help: "The total number of delete requests to the store",
	})
)
