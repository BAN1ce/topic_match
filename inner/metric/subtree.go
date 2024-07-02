package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Write metric

var (
	SubTreeWriteHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "subtree_write_histogram",
		Help: "subtree write histogram",
	})

	SubTreeWriteFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "subtree_write_failed",
		Help: "subtree write failed",
	})
)

// Read metric

var (
	SubTreeReadHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "subtree_read_histogram",
		Help: "subtree read histogram",
	})

	SubTreeReadFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "subtree_read_failed",
		Help: "subtree read failed",
	})

	SubTreeReadNoSubscriber = promauto.NewCounter(prometheus.CounterOpts{
		Name: "subtree_read_no_subscriber",
		Help: "subtree read no subscriber",
	})
)

// Delete metric

var (
	SubTreeDeleteHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "subtree_delete_histogram",
		Help: "subtree delete histogram",
	})

	SubTreeDeleteFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "subtree_delete_failed",
		Help: "subtree delete failed",
	})
)
