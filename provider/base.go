package provider

import (
	"context"

	"github.com/yourlogarithm/l337/internal/chat"
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
	Chat(ctx context.Context, request *chat.Request, options *ChatOptions) (chat.Response, error)
}
