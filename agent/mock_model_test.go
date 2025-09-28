package agent_test

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

type MockModel struct {
	ChatFunc          func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error)
	ChatStreamingFunc func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.ResponseChannel, error)
}

func (m MockModel) Chat(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
	return m.ChatFunc(ctx, req, opts)
}

func (m MockModel) ChatStreaming(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.ResponseChannel, error) {
	return m.ChatStreamingFunc(ctx, req, opts)
}

func (m MockModel) Wrap() *provider.Model {
	return &provider.Model{
		Name:     "mockName",
		Provider: "mockProvider",
		Impl:     m,
	}
}
