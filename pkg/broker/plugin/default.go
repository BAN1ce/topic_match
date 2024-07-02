package plugin

import (
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/eclipse/paho.golang/packets"
)

var (
	defaultOnReceivedConnect = []OnReceivedConnect{
		func(clientID string, connect *packets.Connect) error {
			metric.ReceivedConnect.Add(1)
			return nil
		},
	}

	defaultOnSendConnAck = []OnSendConnAck{
		func(clientID string, connAck *packets.Connack) error {
			metric.SendConnectAck.Add(1)
			return nil
		},
	}

	defaultOnReceivedAuth = []OnReceivedAuth{
		func(clientID string, auth *packets.Auth) error {
			metric.ReceivedAuth.Add(1)
			return nil
		},
	}

	defaultOnReceivedDisconnect = []OnReceivedDisconnect{
		func(clientID string, disconnect *packets.Disconnect) error {
			metric.ReceivedDisconnect.Add(1)
			return nil
		},
	}

	defaultOnReceivedSubscribe = []OnSubscribe{
		func(clientID string, subscribe *packets.Subscribe) error {
			metric.ReceivedSubscription.Add(1)
			return nil
		},
	}

	defaultOnSendSubAck = []OnSendSubAck{
		func(clientID string, subAck *packets.Suback) error {
			metric.SendSubscriptionAck.Add(1)
			return nil
		},
	}

	defaultOnReceivedUnsubscribe = []OnUnsubscribe{
		func(clientID string, unsubscribe *packets.Unsubscribe) error {
			metric.ReceivedUnsubscription.Add(1)
			return nil
		},
	}

	defaultOnSendUnsubAck = []OnSendUnsubAck{
		func(clientID string, unsubAck *packets.Unsuback) error {
			metric.SendUnsubscriptionAck.Add(1)
			return nil
		},
	}

	defaultOnReceivedPublish = []OnReceivedPublish{
		func(clientID string, publish packets.Publish) error {
			metric.ReceivedPublish.Add(1)
			return nil
		},
	}

	defaultOnReceivedPubAck = []OnReceivedPubAck{
		func(clientID string, pubAck packets.Puback) error {
			metric.ReceivedPublishAck.Add(1)
			return nil
		},
	}

	defaultOnReceivedPubRel = []OnReceivedPubRel{
		func(clientID string, pubRel packets.Pubrel) error {
			metric.ReceivedPubRel.Add(1)
			return nil
		},
	}

	defaultOnReceivedPubRec = []OnReceivedPubRec{
		func(clientID string, pubRec packets.Pubrec) error {
			metric.ReceivedPubRec.Add(1)
			return nil
		},
	}

	defaultOnReceivedPubComp = []OnReceivedPubComp{
		func(clientID string, pubComp packets.Pubcomp) error {
			metric.ReceivedPubComp.Add(1)
			return nil
		},
	}

	defaultOnSendPublish = []OnSendPublish{
		func(clientID string, publish *packets.Publish) error {
			metric.SendPublish.Add(1)
			return nil
		},
	}

	defaultOnSendPubAck = []OnSendPubAck{
		func(clientID string, pubAck *packets.Puback) error {
			metric.SendPublishAck.Add(1)
			return nil
		},
	}

	defaultOnSendPubRel = []OnSendPubRel{
		func(clientID string, pubRel *packets.Pubrel) error {
			metric.SendPubRel.Add(1)
			return nil
		},
	}

	defaultOnSendPubRec = []OnSendPubRec{
		func(clientID string, pubRec *packets.Pubrec) error {
			metric.SendPubRec.Add(1)
			return nil
		},
	}

	defaultOnSendPubComp = []OnSendPubComp{
		func(clientID string, pubComp *packets.Pubcomp) error {
			metric.SendPubComp.Add(1)
			return nil
		},
	}

	defaultOnReceivedPingReq = []OnReceivedPingReq{
		func(clientID string, pingReq *packets.Pingreq) error {
			metric.ReceivedPing.Add(1)
			return nil
		},
	}

	defaultOnSendPingResp = []OnSendPingResp{
		func(clientID string, pingResp *packets.Pingresp) error {
			metric.SendPong.Add(1)
			return nil
		},
	}
)

func NewDefaultPlugin() *Plugins {
	return &Plugins{
		PacketPlugin: PacketPlugin{
			OnReceivedConnect: defaultOnReceivedConnect,
			OnSendConnAck:     defaultOnSendConnAck,

			OnReceivedAuth:       defaultOnReceivedAuth,
			OnReceivedDisconnect: defaultOnReceivedDisconnect,

			// about subscribe
			OnSubscribe:  defaultOnReceivedSubscribe,
			OnSendSubAck: defaultOnSendSubAck,

			// about unsubscribe
			OnUnsubscribe:  defaultOnReceivedUnsubscribe,
			OnSendUnsubAck: defaultOnSendUnsubAck,

			// about publish
			OnReceivedPublish: defaultOnReceivedPublish,
			OnReceivedPubAck:  defaultOnReceivedPubAck,

			OnReceivedPubRel:  defaultOnReceivedPubRel,
			OnReceivedPubRec:  defaultOnReceivedPubRec,
			OnReceivedPubComp: defaultOnReceivedPubComp,

			OnSendPublish: defaultOnSendPublish,
			OnSendPubAck:  defaultOnSendPubAck,
			OnSendPubRel:  defaultOnSendPubRel,
			OnSendPubRec:  defaultOnSendPubRec,
			OnSendPubComp: defaultOnSendPubComp,

			// about ping pong
			OnReceivedPingReq: defaultOnReceivedPingReq,
			OnSendPingResp:    defaultOnSendPingResp,
		},
		ClientPlugin: ClientPlugin{},
	}

}
