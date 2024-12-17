package driver

import "github.com/kataras/go-events"

type Local struct {
	e events.EventEmmiter
}

func NewLocal() *Local {
	return &Local{e: events.New()}
}

func (l *Local) AddListener(eventName string, listener events.Listener) {
	l.e.AddListener(events.EventName(eventName), listener)
}

func (l *Local) AddListenerOnce(eventName string, listener events.Listener) {
	l.e.Once(events.EventName(eventName), listener)
}

func (l *Local) RemoveListener(eventName string, listener events.Listener) {
	l.e.RemoveListener(events.EventName(eventName), listener)
}

func (l *Local) Emit(eventName string, args ...interface{}) {
	l.e.Emit(events.EventName(eventName), args...)
}
