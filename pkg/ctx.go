package pkg

import (
	"context"
	"github.com/google/uuid"
)

var (
	ContextIDKey = struct {
	}{}
	ClientIDKey = ClientID{}
)

type ClientID struct {
}

type ContextID struct {
}

func NewCtxWithID() context.Context {
	return context.WithValue(context.Background(), ContextIDKey, uuid.NewString())
}

func GetContextID(ctx context.Context) (id string) {
	if id, ok := ctx.Value(ContextIDKey).(string); ok {
		return id
	}
	return ""
}

func GetClientUID(ctx context.Context) (id string) {
	if id, ok := ctx.Value(ClientIDKey).(string); ok {
		return id
	}
	return ""
}
