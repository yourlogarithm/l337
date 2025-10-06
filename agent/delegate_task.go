package agent

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/tools"
)

const DELEGATE_TASK_TOOL_NAME = "delegate_task"
const DELEGATE_TASK_TOOL_DESC = "Delegates the task to one or more subordinates"

type delegateTaskParams struct {
	// Names of the subordinates to delegate the task to
	Names []string `json:"names" jsonschema:"required,description=Names of the subordinates to delegate the task to"`
	// Expected output from the subordinates
	ExpectedOutput string `json:"expected_output" jsonschema:"required,description=Expected output from the subordinates"`
}

func addDelegateTaskTool(agent *Agent) error {
	if agent == nil {
		return ErrBuilderParams{
			Param: "agent",
			Msg:   "agent is required",
		}
	}

	delegateTask := func(ctx context.Context, response *chat.RunResponse, delegateTaskParams delegateTaskParams) ([]chat.Content, error) {
		logger.Debug(DELEGATE_TASK_TOOL_NAME, "params", delegateTaskParams)

		if len(delegateTaskParams.Names) == 0 {
			return nil, fmt.Errorf("no subordinates specified")
		}

		if delegateTaskParams.ExpectedOutput == "" {
			return nil, fmt.Errorf("no expected output specified")
		}

		nameSet := make(map[string]struct{}, len(delegateTaskParams.Names))
		for _, name := range delegateTaskParams.Names {
			nameSet[name] = struct{}{}
		}

		msg := chat.Message{
			Role:    chat.RoleUser,
			Content: chat.NewTextContent(delegateTaskParams.ExpectedOutput),
		}

		var wg sync.WaitGroup

		responses := make([]chat.Content, 0, len(delegateTaskParams.Names))

		for i := range agent.subordinates {
			name, err := agent.subordinates[i].Name()
			if err != nil {
				return nil, err
			}
			if _, exists := nameSet[name]; exists {
				wg.Add(1)
				go func(sub AgentImpl, name string) {
					defer wg.Done()
					subordinateRunResponse := &chat.RunResponse{
						SessionID: response.SessionID,
						Messages:  []chat.Message{msg},
						Metrics:   make(map[uuid.UUID][]metrics.Metrics),
					}
					err := sub.Run(ctx, subordinateRunResponse)
					for id, metrics := range subordinateRunResponse.Metrics {
						if v, ok := response.Metrics[id]; ok {
							response.Metrics[id] = append(v, metrics...)
						} else {
							response.Metrics[id] = metrics
						}
					}
					if err != nil {
						responses = append(responses, chat.NewTextContent(fmt.Sprintf("(%s) Error: %s", name, err.Error())))
					} else {
						responses = append(responses, subordinateRunResponse.Content())
					}
				}(agent.subordinates[i], name)
			}
		}

		wg.Wait()

		return responses, nil
	}

	tool, err := tools.NewWithArgs(DELEGATE_TASK_TOOL_NAME, DELEGATE_TASK_TOOL_DESC, delegateTask)
	if err != nil {
		panic(err)
	}

	agent.tools.AddTool(tool)

	return nil
}
