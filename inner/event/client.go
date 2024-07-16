package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
)

const (
	ClientConnect = "event.client.connect.request"

	ClientCreateSuccess = "event.client.create.success"
	ClientCreateFailed  = "event.client.create.failed"

	ClientDeleteSuccess = "event.client.delete.success"

	ClientOCountState = "event.client.count.inner"
)

var (
	ClientEvent = newClientEvent()
)

type Client struct {
}

func newClientEvent() *Client {
	c := &Client{}

	return c
}

func (c *Client) EmitClientConnectResult(uid string, ok bool) {

	eventDriver.Emit(ClientConnect)
	if ok {
		eventDriver.Emit(ClientCreateSuccess, uid)
	} else {
		eventDriver.Emit(ClientCreateFailed, uid)
	}
}

func (c *Client) EmitClientCountState(count int64) {
	eventDriver.Emit(ClientOCountState, count)
}

func (c *Client) EmitClientDeleteSuccess(uid string) {
	eventDriver.Emit(ClientDeleteSuccess, uid)

}

func (c *Client) CreateListenClientOnline(handler func(i ...interface{})) {
	eventDriver.AddListener(ClientCreateSuccess, handler)
}

func (c *Client) CreateListenClientOffline(handler func(i ...interface{})) {
	eventDriver.AddListener(ClientDeleteSuccess, handler)
}

func (c *Client) addDefaultListener() {
	c.CreateListenClientOnline(func(i ...interface{}) {
		if len(i) < 1 {
			return
		}
		if _, ok := i[0].(string); ok {
			metric.ClientConnectSuccessEvent.Inc()
		}
	})

	c.CreateListenClientOffline(func(i ...interface{}) {
		if len(i) < 1 {
			return
		}
		if _, ok := i[0].(string); ok {
			metric.ClientOnline.Inc()
		}
	})
}
