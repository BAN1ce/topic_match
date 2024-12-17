package broker

import (
	"context"
	"github.com/BAN1ce/skyTree/pkg/packet"
)

type MessageSource interface {
	ReadMessageSource
	StreamSource
}

type ReadMessageSource interface {
	NextMessages(ctx context.Context, n int, startMessageID string, include bool, in chan *packet.Message) ([]*packet.Message, error)
}

type StreamSource interface {
	ListenMessage(ctx context.Context, writer chan *packet.Message) error
	Close() error
}
