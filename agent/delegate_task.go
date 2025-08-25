package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

type delegateTaskParams struct {
	// Names of the subordinates to delegate the task to
	Names []string `json:"names" jsonschema:"required,description=Names of the subordinates to delegate the task to"`
	// Expected output from the subordinates
	ExpectedOutput string `json:"expected_output" jsonschema:"required,description=Expected output from the subordinates"`
}

func addDelegateTaskTool(agent *Agent) error {
	const delegate_task_tool_name = "delegate_task"

	if agent == nil {
		return fmt.Errorf("agent is nil")
	}

	delegateTask := func(ctx context.Context, response *run.Response, delegateTaskParams delegateTaskParams) (string, error) {
		logger.Debug(delegate_task_tool_name, "params", delegateTaskParams)

		if len(delegateTaskParams.Names) == 0 {
			return "", fmt.Errorf("no subordinates specified")
		}

		if delegateTaskParams.ExpectedOutput == "" {
			return "", fmt.Errorf("no expected output specified")
		}

		nameSet := make(map[string]struct{}, len(delegateTaskParams.Names))
		for _, name := range delegateTaskParams.Names {
			nameSet[name] = struct{}{}
		}

		msg := chat.Message{
			Role:    chat.RoleUser,
			Content: delegateTaskParams.ExpectedOutput,
		}

		var wg sync.WaitGroup

		var sb strings.Builder

		for i := range agent.subordinates {
			if _, exists := nameSet[agent.subordinates[i].Name()]; exists {
				wg.Add(1)
				go func(sub AgentImpl) {
					defer wg.Done()
					subordinateRunResponse := &run.Response{
						SessionID: response.SessionID,
						Messages:  []chat.Message{msg},
						Metrics:   make(map[uuid.UUID][]metrics.Metrics),
					}
					err := sub.run(ctx, subordinateRunResponse)
					for id, metrics := range subordinateRunResponse.Metrics {
						if v, ok := response.Metrics[id]; ok {
							response.Metrics[id] = append(v, metrics...)
						} else {
							response.Metrics[id] = metrics
						}
					}
					if err != nil {
						sb.WriteString(fmt.Sprintf("(%s) Error: %s\n", sub.Name(), err.Error()))
					} else {
						sb.WriteString(fmt.Sprintf("(%s) Response: %s\n", sub.Name(), subordinateRunResponse.Content()))
					}
				}(agent.subordinates[i])
			}
		}

		wg.Wait()

		return sb.String(), nil
	}

	tool, err := tools.NewToolWithArgs(delegate_task_tool_name, "Delegates the task to one or more subordinates", delegateTask)
	if err != nil {
		panic(err)
	}

	agent.tools.AddTool(tool)

	return nil
}
