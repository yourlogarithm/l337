package provider

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
)

type Model struct {
	// LLM identifier
	Name string
	// Provider name (e.g., "ollama", "openai")
	Provider string
	// Model implementation
	Impl ModelImpl
}

type ModelImpl interface {
	Chat(context.Context, *Request, *chat.Options) (Response, error)
	ChatStreaming(context.Context, *Request, *chat.Options) (ResponseChannel, error)
}
