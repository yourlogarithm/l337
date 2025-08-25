package agent

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	internal_chat "github.com/yourlogarithm/l337/internal/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/retry"
	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

// Convenience method to run an agent without [run.Parameters] declaration
func (a *Agent) RunWithParams(ctx context.Context, params ...run.Parameter) (run.Response, error) {
	var runParams run.Parameters
	for _, param := range params {
		if err := param.Apply(&runParams); err != nil {
			return run.Response{}, err
		}
	}
	if len(runParams.Messages) == 0 {
		return run.Response{}, fmt.Errorf("no messages provided")
	}

	runResponse := &run.Response{
		SessionID: runParams.SessionID,
		Messages:  runParams.Messages,
		Metrics:   make(map[uuid.UUID][]metrics.Metrics),
	}

	return *runResponse, a.run(ctx, runResponse)
}

func (a *Agent) run(ctx context.Context, runResponse *run.Response) error {
	if a.retry == nil {
		a.retry = retry.Default()
	}

	if len(runResponse.Messages) == 0 || runResponse.Messages[0].Role != chat.RoleSystem {
		systemMsg := chat.Message{
			Role:    chat.RoleSystem,
			Content: a.ComputeSystemMessage(),
		}
		runResponse.Messages = slices.Insert(runResponse.Messages, 0, systemMsg)
	}

	tools := make([]tools.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, tool)
	}

	for {
		var chatResponse internal_chat.Response
		req := internal_chat.Request{
			Messages: runResponse.Messages,
			Tools:    tools,
		}
		logger.Debug("agent.run.request", "agent", a.name, "request", req)
		if err := a.retry.Execute(func() error {
			response, err := a.model.Impl.Chat(ctx, &req, &a.chatOptions)
			if err != nil {
				return err
			}
			chatResponse = response
			return nil
		}); err != nil {
			return err
		}
		logger.Debug("agent.run.response", "agent", a.name, "response", chatResponse)
		msg := chat.Message{
			Role:      chat.RoleAssistant,
			Content:   chatResponse.Content,
			ToolCalls: chatResponse.ToolCalls,
		}
		runResponse.Messages = append(runResponse.Messages, msg)
		chatResponse.Metrics.SessionID = runResponse.SessionID
		if v, ok := runResponse.Metrics[a.id]; ok {
			runResponse.Metrics[a.id] = append(v, chatResponse.Metrics)
		} else {
			runResponse.Metrics[a.id] = []metrics.Metrics{chatResponse.Metrics}
		}
		if len(chatResponse.ToolCalls) > 0 {
			var wg sync.WaitGroup
			var mu sync.Mutex

			wg.Add(len(chatResponse.ToolCalls))

			type ToolCallResult struct {
				ToolCall *chat.ToolCall
				Content  string
				IsErr    bool
			}

			results := make(map[string]ToolCallResult, len(chatResponse.ToolCalls))
			order := make([]string, 0, len(chatResponse.ToolCalls))

			for _, toolCall := range chatResponse.ToolCalls {
				order = append(order, toolCall.ID)
				tc := toolCall
				go func(toolCall *chat.ToolCall) {
					defer wg.Done()
					var content string
					var isErr bool

					tool, exists := a.tools.Get(toolCall.Name)
					if exists {
						result, err := tool.Callable(ctx, runResponse, toolCall.Arguments)
						if err != nil {
							content = "error: " + err.Error()
							isErr = true
						} else {
							content = result
						}
					} else {
						content = fmt.Sprintf("error: tool not found: %s", toolCall.Name)
					}
					mu.Lock()
					defer mu.Unlock()
					results[toolCall.ID] = ToolCallResult{
						ToolCall: toolCall,
						Content:  content,
						IsErr:    isErr,
					}
				}(&tc)
			}
			wg.Wait()

			for _, id := range order {
				result, exists := results[id]
				if !exists {
					return fmt.Errorf("tool call result not found for ID: %s", id)
				}
				runResponse.Messages = append(runResponse.Messages, chat.Message{
					Role:    chat.RoleTool,
					Content: result.Content,
					Name:    result.ToolCall.Name,
					IsErr:   result.IsErr,
				})
			}
		} else {
			break
		}
	}

	return nil
}
