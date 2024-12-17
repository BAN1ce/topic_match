package event

import (
	"github.com/kataras/go-events"
	"sync"
)

type Driver interface {
	AddListener(eventName string, listener events.Listener)
	AddListenerOnce(eventName string, listener events.Listener)
	RemoveListener(eventName string, listener events.Listener)
	Emit(eventName string, args ...interface{})
}

type defaultListener interface {
	addDefaultListener()
}

var (
	eventDriver     *wrapperDriver
	mux             sync.RWMutex
	registerMetrics = []defaultListener{
		MQTTEvent,
		StoreEvent,
		SubtreeEvent,
		ClientEvent,
	}
)

func SetEventDriverOnce(driver Driver) {
	events.New()
	mux.Lock()
	defer mux.Unlock()
	if eventDriver == nil {
		eventDriver = &wrapperDriver{driver: driver}
		for _, metric := range registerMetrics {
			metric.addDefaultListener()
		}
	}
}

type wrapperDriver struct {
	driver Driver
}

func (w *wrapperDriver) AddListener(eventName string, listener events.Listener) {
	w.driver.AddListener(eventName, listener)
}

func (w *wrapperDriver) RemoveListener(eventName string, listener events.Listener) {
	w.driver.RemoveListener(eventName, listener)
}

func (w *wrapperDriver) Emit(eventName string, args ...interface{}) {
	w.driver.Emit(eventName, args...)
}

func (w *wrapperDriver) AddListenerOnce(eventName string, listener events.Listener) {
	w.driver.AddListenerOnce(eventName, listener)
}
