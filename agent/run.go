package agent

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

// Convenience method to run an agent without [run.Parameters] declaration
func (a *Agent) RunWithParams(ctx context.Context, params ...types.Parameter) (*types.Run, error) {
	run, err := BuildRun(params...)
	if err != nil {
		return nil, err
	}
	return run, a.Run(ctx, run)
}

func (a *Agent) Run(ctx context.Context, run *types.Run) error {
	tools, err := a.preRun(run)
	if err != nil {
		return err
	}

	for {
		logger.Debug("agent.run.request", "agent", a.name, "messages", run.Messages, "tools", tools)
		chatResponse, err := a.model.Impl.Chat(ctx, run.Messages, tools, &a.chatOptions)
		if err != nil {
			return err
		}
		toolsCalled, err := a.handleResponse(ctx, run, &chatResponse)
		if err != nil {
			return err
		}
		if !toolsCalled {
			break
		}
	}

	return nil
}

// preRun performs common pre-run checks and setups.
func (a *Agent) preRun(run *types.Run) ([]tools.Tool, error) {
	if a.model == nil {
		return nil, ErrBuilderParams{
			Param: "model",
			Msg:   "nil",
		}
	}

	if len(run.Messages) == 0 {
		return nil, ErrBuilderParams{
			Param: "messages",
			Msg:   "at least one message is required",
		}
	}

	if err := a.addSystemMessageIfMissing(run); err != nil {
		return nil, err
	}

	return a.Tools()
}

func (a *Agent) addSystemMessageIfMissing(run *types.Run) error {
	if len(run.Messages) == 0 || run.Messages[0].Role != types.RoleSystem {
		logger.Debug("adding system message to run response")
		content, err := a.ComputeSystemMessage()
		if err != nil {
			return err
		}
		if content != "" {
			systemMsg := types.Message{
				Role:    types.RoleSystem,
				Content: types.NewTextContent(content),
			}
			run.Messages = slices.Insert(run.Messages, 0, systemMsg)
		}
	}
	return nil
}

func (a *Agent) handleResponse(ctx context.Context, run *types.Run, chatResponse *types.Response) (bool, error) {
	logger.Debug("agent.run.response", "agent", a.name, "response", chatResponse)

	if chatResponse == nil {
		return false, ErrModelResponse{
			Msg: "nil",
		}
	}

	if chatResponse.FinishReason == "" {
		return false, ErrModelResponse{
			Msg: "no finish reason",
		}
	}

	msg := types.Message{
		Role:      types.RoleAssistant,
		Content:   chatResponse.Content,
		Reasoning: chatResponse.Reasoning,
		ToolCalls: chatResponse.ToolCalls,
	}
	run.Messages = append(run.Messages, msg)
	chatResponse.Metrics.SessionID = run.SessionID

	if run.Metrics == nil {
		run.Metrics = make(map[uuid.UUID][]metrics.Metrics)
	}
	if v, ok := run.Metrics[a.id]; ok {
		run.Metrics[a.id] = append(v, chatResponse.Metrics)
	} else {
		run.Metrics[a.id] = []metrics.Metrics{chatResponse.Metrics}
	}

	if len(chatResponse.ToolCalls) > 0 {
		var wg sync.WaitGroup
		var mu sync.Mutex

		wg.Add(len(chatResponse.ToolCalls))

		type ToolCallResult struct {
			ToolCall *types.ToolCall
			Content  types.Content
			IsErr    bool
		}

		tool_call_results_map := make(map[string][]ToolCallResult, len(chatResponse.ToolCalls))
		order := make([]string, 0, len(chatResponse.ToolCalls))

		for _, toolCall := range chatResponse.ToolCalls {
			order = append(order, toolCall.ID)
			tc := toolCall
			go func(toolCall *types.ToolCall) {
				defer wg.Done()
				var contents []types.Content
				var isErr bool

				tool, exists := a.tools.Get(toolCall.Name)
				if exists {
					result, err := tool.Call(ctx, run, toolCall.Arguments)
					if err != nil {
						contents = append(contents, types.NewTextContent("error: "+err.Error()))
						isErr = true
					} else {
						contents = append(contents, result...)
					}
				} else {
					contents = append(contents, types.NewTextContent(fmt.Sprintf("error: tool not found: %s", toolCall.Name)))
				}
				mu.Lock()
				defer mu.Unlock()
				tool_call_results_map[toolCall.ID] = make([]ToolCallResult, 0, len(contents))
				for _, content := range contents {
					tool_call_results_map[toolCall.ID] = append(tool_call_results_map[toolCall.ID], ToolCallResult{
						ToolCall: toolCall,
						Content:  content,
						IsErr:    isErr,
					})
				}
			}(&tc)
		}
		wg.Wait()

		for _, id := range order {
			tool_call_results, exists := tool_call_results_map[id]
			if !exists {
				return false, ErrModelResponse{
					Msg: fmt.Sprintf("tool call result not found for ID: %s", id),
				}
			}
			for _, result := range tool_call_results {
				run.Messages = append(run.Messages, types.Message{
					Role:    types.RoleTool,
					Content: result.Content,
					Name:    result.ToolCall.ID,
					IsErr:   result.IsErr,
				})
			}
		}

		return true, nil
	} else {
		return false, nil
	}
}
