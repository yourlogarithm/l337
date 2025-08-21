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
	// Names of the subordinates to delegate the task to
	Names []string `json:"names" jsonschema:"required,description=Names of the subordinates to delegate the task to"`
	// Expected output from the subordinates
	ExpectedOutput string `json:"expected_output" jsonschema:"required,description=Expected output from the subordinates"`
}

func addDelegateTaskTool(agent *Agent) {
	delegateTask := func(ctx context.Context, params delegateTaskParams) (string, error) {
		logger.Debug("delegate_task", "params", params)

		if len(params.Names) == 0 {
			return "", fmt.Errorf("no subordinates specified")
		}

		if params.ExpectedOutput == "" {
			return "", fmt.Errorf("no expected output specified")
		}

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

	tool, err := tools.NewToolWithArgs("delegate_task", "Delegates the task to one or more subordinates", delegateTask)
	if err != nil {
		panic(err)
	}

	agent.tools.AddTool(tool)
}
