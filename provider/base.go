package provider

import (
	"context"

	"github.com/yourlogarithm/golagno/chat"
)

type Model struct {
	Name     string
	Provider string
	Impl     ModelImpl
}

type ModelImpl interface {
	Chat(ctx context.Context, messages []chat.Message) (chat.Response, error)
}
