package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ClientOnline = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_online_state",
	})

	ClientConnectEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_connect_event",
	})

	ClientConnectSuccessEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_connect_success_event",
	})

	ClientConnectFailedEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_connect_failed_event",
	})

	ClientDisconnectRequestEvent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "client_disconnect_event",
	})
)
