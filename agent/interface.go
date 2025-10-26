package agent

import (
	"context"

	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

type AgentImpl interface {
	Name() (string, error)
	Description() (string, error)
	Tools() ([]tools.Tool, error)
	Run(context.Context, *types.Run) error
	RunStreaming(context.Context, *types.Run, int) (streaming.ResponseChannel, error)
}
