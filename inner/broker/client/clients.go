package client

import (
	"github.com/BAN1ce/skyTree/inner/metric"
	"github.com/BAN1ce/skyTree/logger"
	"github.com/BAN1ce/skyTree/pkg/retry"
	"sync"
)

type ID = string

type Clients struct {
	mux     sync.RWMutex
	clients map[ID]*Client
}

func NewClients() *Clients {
	return &Clients{
		clients: map[ID]*Client{},
	}
}

// CreateClient creates a new client and returns the client
// if the client already exists,
// it will return the old client and return false
// if the client is created successfully,
// it will return the new client and return true
func (c *Clients) CreateClient(client *Client) (*Client, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return c.createClient(client)
}

func (c *Clients) DestroyClient(id ID) (*Client, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if client2, ok := c.clients[id]; ok {
		logger.Logger.Debug().Str("id", id).Msg("delete client")
		delete(c.clients, id)
		if err := client2.Close(); err != nil {
			logger.Logger.Warn().Err(err).Str("id", id).Msg("close client error")
		}
		return client2, ok
	}

	logger.Logger.Warn().Str("id", id).Msg("delete client not found")
	return nil, false
}

func (c *Clients) createClient(client *Client) (*Client, bool) {
	if oldClient, ok := c.clients[client.GetID()]; ok {
		if oldClient != client {
			logger.Logger.Warn().Str("client uid", client.GetID()).Msg("client already exist")
			if err := oldClient.Close(); err != nil {
				logger.Logger.Warn().Err(err).Str("uid", client.GetID()).Msg("close old client error")
			}
			return oldClient, false
		}
	}
	c.clients[client.GetID()] = client
	return client, true
}

func (c *Clients) ReadClient(uid string) (*Client, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	client2, ok := c.clients[uid]
	return client2, ok
}

func (c *Clients) Count() int64 {
	c.mux.RLock()
	tmp := int64(len(c.clients))
	c.mux.RUnlock()

	return tmp
}

func (c *Clients) RetryPublish(task *retry.Task) error {
	if cl, ok := c.ReadClient(task.Data.SendClientID); ok {
		return cl.RetrySend(task.Data)
	}
	return nil

}

func (c *Clients) RetryPublishTimeout(task *retry.Task) error {
	metric.PublishRetryTimeout.Inc()
	metric.PublishRetryTaskCurrent.Add(-1)

	logger.Logger.Debug().Str("clientID", task.ClientID).Msg("retry publish timeout, close client")
	c.DestroyClient(task.ClientID)
	return nil

}
