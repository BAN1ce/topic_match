package client

import (
	"github.com/BAN1ce/skyTree/inner/broker/monitor"
	"github.com/BAN1ce/skyTree/inner/broker/store/message"
	"github.com/BAN1ce/skyTree/inner/broker/will"
	"github.com/BAN1ce/skyTree/inner/event"
	"github.com/BAN1ce/skyTree/pkg/broker"
	"github.com/BAN1ce/skyTree/pkg/broker/plugin"
	"github.com/BAN1ce/skyTree/pkg/broker/retain"
	"github.com/BAN1ce/skyTree/pkg/broker/session"
	"time"
)

type NotifyClientClose interface {
	NotifyClientClose(c *Client)
	NotifyWillMessage(message *session.WillMessage)
}

type ComponentOption func(*Component)

type Component struct {
	willMonitor         will.MessageMonitor
	Store               broker.TopicMessageStore
	session             session.Session
	retain              retain.Retain
	cfg                 Config
	notifyClose         NotifyClientClose
	plugin              *plugin.Plugins
	subCenter           broker.SubCenter
	messageStoreWrapper *message.Wrapper
	monitorKeepAlive    *monitor.KeepAlive
	event               event.Driver
}

func WithEvent(driver event.Driver) ComponentOption {
	return func(options *Component) {
		options.event = driver
	}

}

func WithMessageStoreWrapper(wrapper *message.Wrapper) ComponentOption {
	return func(component *Component) {
		component.messageStoreWrapper = wrapper
	}
}

func WithRetain(retain2 retain.Retain) ComponentOption {
	return func(options *Component) {
		options.retain = retain2
	}
}

func WithStore(store broker.TopicMessageStore) ComponentOption {
	return func(options *Component) {
		options.Store = store
	}
}
func WithConfig(cfg Config) ComponentOption {
	return func(options *Component) {
		options.cfg = cfg
	}
}

func WithNotifyClose(notifyClose NotifyClientClose) ComponentOption {
	return func(options *Component) {
		options.notifyClose = notifyClose
	}
}

func WithPlugin(plugins *plugin.Plugins) ComponentOption {
	return func(options *Component) {
		options.plugin = plugins
	}
}

func WithSession(session2 session.Session) ComponentOption {
	return func(options *Component) {
		options.session = session2
	}
}

func WithKeepAliveTime(keepAlive time.Duration) ComponentOption {
	return func(options *Component) {
		options.cfg.KeepAlive = keepAlive
	}
}

func WithSubCenter(subCenter broker.SubCenter) ComponentOption {
	return func(c *Component) {
		c.subCenter = subCenter
	}
}

func WithMonitorKeepAlive(monitorKeepAlive *monitor.KeepAlive) ComponentOption {
	return func(c *Component) {
		c.monitorKeepAlive = monitorKeepAlive
	}
}
