package providers

import (
	"context"

	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
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
	Chat(context.Context, []types.Message, []tools.Tool, *types.Options) (types.Response, error)
	ChatStreaming(context.Context, []types.Message, []tools.Tool, *types.Options) (streaming.ResponseChannel, error)
}
