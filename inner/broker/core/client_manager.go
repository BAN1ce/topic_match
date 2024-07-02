package core

import (
	"github.com/BAN1ce/skyTree/inner/broker/client"
)

type UID = string

func NewClientManager() *client.Clients {
	return client.NewClients()
}
