package agent

import (
	"context"

	"github.com/yourlogarithm/golagno/agentic"
	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/run"
)

func (a *Agent) Run(ctx context.Context, messages []chat.Message) (run.Response, error) {
	newMessages := []chat.Message{
		{
			Role:    chat.RoleSystem.String(),
			Content: a.Options.ComputeSystemMessage(),
		},
	}

	newMessages = append(newMessages, messages...)

	return agentic.Run(ctx, newMessages, &a.Options, logger)
}
