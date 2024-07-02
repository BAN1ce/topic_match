package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Connect packet

var (
	ReceivedConnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_connect",
	})

	SendConnectAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_connect_ack",
	})
)

// Subscribe Packet

var (
	ReceivedSubscription = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_subscription",
	})

	SendSubscriptionAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_subscription_ack",
	})
)

// Unsubscribe Packet

var (
	ReceivedUnsubscription = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_unsubscription",
	})

	SendUnsubscriptionAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_unsubscription_ack",
	})
)

// Publish Packet

var (
	ReceivedPublish = promauto.NewCounter(prometheus.CounterOpts{
		Name: "receive_publish",
	})

	SendPublish = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_publish",
	})
)

// Publish Ack

var (
	SendPublishAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_publish_ack",
	})

	ReceivedPublishAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_publish_ack",
	})
)

// Publish Rec
var (
	ReceivedPubRec = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_pubrec",
	})

	SendPubRec = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_pubrec",
	})
)

// Publish Rel
var (
	ReceivedPubRel = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_pubrel",
	})

	SendPubRel = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_pubrel",
	})
)

// Publish Comp

var (
	ReceivedPubComp = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_pubcomp",
	})

	SendPubComp = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_pubcomp",
	})
)

var (
	ReceivedAuth = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_auth",
	})

	SendAuthAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_auth_ack",
	})
)

// Disconnect

var (
	ReceivedDisconnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_disconnect",
	})

	SendDisconnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_disconnect",
	})
)

// Ping

var (
	ReceivedPing = promauto.NewCounter(prometheus.CounterOpts{
		Name: "received_ping",
	})

	SendPong = promauto.NewCounter(prometheus.CounterOpts{
		Name: "send_pong",
	})
)
