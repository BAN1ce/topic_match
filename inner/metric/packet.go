package metric

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Connect packet

var (
	ReceivedConnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_connect",
	})

	SendConnectAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_connect_ack",
	})
)

// Subscribe Packet

var (
	ReceivedSubscription = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_subscription",
	})

	SendSubscriptionAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_subscription_ack",
	})
)

// Unsubscribe Packet

var (
	ReceivedUnsubscription = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_unsubscription",
	})

	SendUnsubscriptionAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_unsubscription_ack",
	})
)

// Publish Packet

var (
	ReceivedPublish = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_receive_publish",
	})

	SendPublish = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_publish",
	})
)

// Publish Ack

var (
	SendPublishAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_publish_ack",
	})

	ReceivedPublishAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_publish_ack",
	})
)

// Publish Rec

var (
	ReceivedPubRec = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_pubrec",
	})

	SendPubRec = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_pubrec",
	})
)

// Publish Rel

var (
	ReceivedPubRel = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_pubrel",
	})

	SendPubRel = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_pubrel",
	})
)

// Publish Comp

var (
	ReceivedPubComp = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_pubcomp",
	})

	SendPubComp = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_pubcomp",
	})
)

// Auth

var (
	ReceivedAuth = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_auth",
	})

	SendAuthAck = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_auth_ack",
	})
)

// Disconnect

var (
	ReceivedDisconnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_disconnect",
	})

	SendDisconnect = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_disconnect",
	})
)

// Ping

var (
	ReceivedPing = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_received_ping",
	})

	SendPong = promauto.NewCounter(prometheus.CounterOpts{
		Name: "mqtt_send_pong",
	})
)
