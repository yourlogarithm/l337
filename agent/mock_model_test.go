package agent_test

import (
	"context"

	"github.com/yourlogarithm/l337/providers"
	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

type MockModel struct {
	ChatFunc          func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error)
	ChatStreamingFunc func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error)
}

func (m MockModel) Chat(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
	return m.ChatFunc(ctx, messages, tools, opts)
}

func (m MockModel) ChatStreaming(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error) {
	return m.ChatStreamingFunc(ctx, messages, tools, opts)
}

func (m MockModel) Wrap() *providers.Model {
	return &providers.Model{
		Name:     "mockName",
		Provider: "mockProvider",
		Impl:     m,
	}
}
