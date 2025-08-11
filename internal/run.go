package internal

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/chat"
	internal_chat "github.com/yourlogarithm/l337/internal/chat"
	"github.com/yourlogarithm/l337/retry"
	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

// Execution loop called internally by the `agent.Agent.Run` or `team.Team.Run` methods.
func Run(ctx context.Context, messages []chat.Message, options *agentic.Options, logger *slog.Logger) (runResponse run.Response, err error) {
	if options == nil {
		return runResponse, fmt.Errorf("options cannot be nil")
	}

	if options.Model == nil {
		return runResponse, fmt.Errorf("model cannot be nil")
	}

	if options.Retry == nil {
		options.Retry = retry.Default()
	}

	runResponse.Messages = append(runResponse.Messages, messages...)

	tools := make([]tools.Tool, 0, len(options.Tools))
	for _, tool := range options.Tools {
		tools = append(tools, tool)
	}

	for {
		var chatResponse internal_chat.Response
		req := internal_chat.Request{
			Messages: runResponse.Messages,
			Tools:    tools,
		}
		logger.Debug("agent.run.request", "agent", options.Name, "request", req)
		if err = options.Retry.Execute(func() error {
			response, err := options.Model.Impl.Chat(ctx, &req)
			if err != nil {
				return err
			}
			chatResponse = response
			return nil
		}); err != nil {
			return runResponse, err
		}
		logger.Debug("agent.run.response", "agent", options.Name, "response", chatResponse)
		msg := chat.Message{
			Role:      chat.RoleAssistant,
			Content:   chatResponse.Content,
			ToolCalls: chatResponse.ToolCalls,
		}
		runResponse.Messages = append(runResponse.Messages, msg)
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

					tool, exists := options.Tools.Get(toolCall.Name)
					if exists {
						result, err := tool.Callable(ctx, toolCall.Arguments)
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
					return runResponse, fmt.Errorf("tool call result not found for ID: %s", id)
				}
				runResponse.AddMessage(chat.Message{
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

	return runResponse, nil
}
