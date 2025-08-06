package team

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/golagno/agentic"
	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/run"
)

func (t *Team) Run(ctx context.Context, messages []chat.Message) (run.Response, error) {
	t.initialize()

	sysMsg, err := t.computeSystemMessage()
	if err != nil {
		return run.Response{}, fmt.Errorf("failed to compute system message: %w", err)
	}

	newMessages := []chat.Message{
		{
			Role:    chat.RoleSystem.String(),
			Content: sysMsg,
		},
	}
	newMessages = append(newMessages, messages...)

	return agentic.Run(ctx, newMessages, t.GetOptions(), logger)
}
