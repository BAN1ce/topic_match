package client

import (
	"github.com/BAN1ce/skyTree/logger"
	"go.uber.org/zap"
	"sync"
)

type UID = string

type Clients struct {
	mux     sync.RWMutex
	clients map[UID]*Client
}

func NewClients() *Clients {
	return &Clients{
		clients: map[string]*Client{},
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

func (c *Clients) DestroyClient(uid string) (*Client, bool) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if client2, ok := c.clients[uid]; ok {
		logger.Logger.Debug("delete client", zap.String("uid", uid))
		delete(c.clients, uid)
		if err := client2.Close(); err != nil {
			logger.Logger.Warn("close client error", zap.Error(err))
		}
		return client2, ok
	}

	logger.Logger.Warn("delete client not found", zap.String("uid", uid))
	return nil, false
}

func (c *Clients) createClient(client *Client) (*Client, bool) {
	if oldClient, ok := c.clients[client.GetUid()]; ok {
		if oldClient != client {
			logger.Logger.Warn("client already exist", zap.String("client uid", client.UID))
			if err := oldClient.Close(); err != nil {
				logger.Logger.Warn("close old client error", zap.Error(err), zap.String("uid", client.UID))
			}
			return oldClient, false
		}
	}
	c.clients[client.GetUid()] = client
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
