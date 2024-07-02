package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Store Read

var (
	StoreReadMessageCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_read_message_count",
		Help: "The total number of messages read from the store",
	})

	StoreReadRequestFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_read_request_failed",
		Help: "The total number of failed read requests to the store",
	})

	StoreReadDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "store_read_duration",
		Help: "The duration of read requests to the store",
		Buckets: []float64{
			0.05, 0.1, 0.5, 1, 2, 5, 10,
		},
	})
)

// Store Save

var (
	StoreStoreRequestFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_store_request_failed",
		Help: "The total number of failed store requests to the store",
	})

	StoreStoreDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "store_store_duration",
		Help: "The duration of store requests to the store",
		Buckets: []float64{
			0.05, 0.1, 0.5, 1, 2, 5, 10,
		},
	})
)

// Store Delete metric

var (
	StoreDeleteRequestFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "store_delete_request_failed",
		Help: "The total number of failed delete requests to the store",
	})

	StoreDeleteDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "store_delete_duration",
		Help: "The duration of delete requests to the store",
		Buckets: []float64{
			0.05, 0.1, 0.5, 1, 2, 5, 10,
		},
	})
)
