package agent

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/internal"
	"github.com/yourlogarithm/l337/run"
)

func (a *Agent) Run(ctx context.Context, messages []chat.Message) (run.Response, error) {
	newMessages := []chat.Message{
		{
			Role:    chat.RoleSystem,
			Content: a.Configuration.ComputeSystemMessage(),
		},
	}

	newMessages = append(newMessages, messages...)

	return internal.Run(ctx, newMessages, &a.Configuration, logger)
}
