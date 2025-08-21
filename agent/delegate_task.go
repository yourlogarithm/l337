package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/tools"
)

type delegateTaskParams struct {
	Names          []string `json:"names"`
	ExpectedOutput string   `json:"expected_output"`
}

func addDelegateTaskTool(agent *Agent) {
	delegateTask := func(ctx context.Context, params delegateTaskParams) (string, error) {
		nameSet := make(map[string]struct{}, len(params.Names))
		for _, name := range params.Names {
			nameSet[name] = struct{}{}
		}

		msg := chat.Message{
			Role:    chat.RoleUser,
			Content: params.ExpectedOutput,
		}
		wrapped := []chat.Message{msg}

		var wg sync.WaitGroup

		var sb strings.Builder

		for i := range agent.subordinates {
			if _, exists := nameSet[agent.subordinates[i].Name()]; exists {
				wg.Add(1)
				go func(sub AgentImpl) {
					defer wg.Done()
					response, err := sub.Run(ctx, wrapped)
					if err != nil {
						sb.WriteString(fmt.Sprintf("(%s) Error: %s\n", sub.Name(), err.Error()))
					} else {
						sb.WriteString(fmt.Sprintf("(%s) Response: %s\n", sub.Name(), response.Content()))
					}
				}(agent.subordinates[i])
			}
		}

		wg.Wait()

		return sb.String(), nil
	}

	tool := tools.NewToolWithArgs("delegate_task", "Delegates the task to one or more subordinates", delegateTask)

	agent.tools.AddTool(tool)
}
