package event

import (
	"github.com/BAN1ce/skyTree/inner/metric"
)

const (
	ClientConnect = "event.client.connect.request"

	ClientCreateSuccess = "event.client.create.success"
	ClientCreateFailed  = "event.client.create.failed"

	ClientDeleteSuccess = "event.client.delete.success"

	ClientOCountState = "event.client.count.state"
)

func (e *Event) EmitClientConnectResult(uid string, ok bool) {

	Driver.Emit(ClientConnect)
	if ok {
		Driver.Emit(ClientCreateSuccess, uid)
	} else {
		Driver.Emit(ClientCreateFailed, uid)
	}
}

func (e *Event) EmitClientCountState(count int64) {
	Driver.Emit(ClientOCountState, count)
}

func (e *Event) EmitClientDeleteSuccess(uid string) {
	Driver.Emit(ClientDeleteSuccess, uid)

}

func (e *Event) CreateListenClientOnline(handler func(i ...interface{})) {
	Driver.AddListener(ClientCreateSuccess, handler)
}

func (e *Event) CreateListenClientOffline(handler func(i ...interface{})) {
	Driver.AddListener(ClientDeleteSuccess, handler)
}

func createClientMetricListen() {

	GlobalEvent.CreateListenClientOnline(func(i ...interface{}) {
		if len(i) < 1 {
			return
		}
		if _, ok := i[0].(string); ok {
			metric.ClientConnectSuccessEvent.Inc()
		}
	})

	GlobalEvent.CreateListenClientOffline(func(i ...interface{}) {
		if len(i) < 1 {
			return
		}
		if _, ok := i[0].(string); ok {
			metric.ClientOnline.Inc()
		}
	})

}
