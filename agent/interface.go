package agent

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

type AgentImpl interface {
	Name() (string, error)
	Description() (string, error)
	Tools() ([]tools.Tool, error)
	Run(context.Context, *chat.RunResponse) error
	RunStreaming(context.Context, *chat.RunResponse, int) (provider.ResponseChannel, error)
}
