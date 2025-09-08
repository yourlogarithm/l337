package agent

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

// Convenience method to run an agent without [run.Parameters] declaration
func (a *Agent) RunWithParams(ctx context.Context, params ...chat.Parameter) (*chat.RunResponse, error) {
	runResponse, err := BuildRunResponse(params...)
	if err != nil {
		return nil, err
	}
	return runResponse, a.Run(ctx, runResponse)
}

func (a *Agent) Run(ctx context.Context, runResponse *chat.RunResponse) error {
	if len(runResponse.Messages) == 0 || runResponse.Messages[0].Role != chat.RoleSystem {
		content, err := a.ComputeSystemMessage()
		if err != nil {
			return err
		}
		systemMsg := chat.Message{
			Role:    chat.RoleSystem,
			Content: content,
		}
		runResponse.Messages = slices.Insert(runResponse.Messages, 0, systemMsg)
	}

	tools := make([]tools.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, tool)
	}

	for {
		req := provider.Request{
			Messages: runResponse.Messages,
			Tools:    tools,
		}
		logger.Debug("agent.run.request", "agent", a.name, "request", req)
		chatResponse, err := a.model.Impl.Chat(ctx, &req, &a.chatOptions)
		if err != nil {
			return err
		}
		toolsCalled, err := a.handleResponse(ctx, runResponse, &chatResponse)
		if err != nil {
			return err
		}
		if !toolsCalled {
			break
		}
	}

	return nil
}

func (a *Agent) handleResponse(ctx context.Context, runResponse *chat.RunResponse, chatResponse *provider.Response) (bool, error) {
	logger.Debug("agent.run.response", "agent", a.name, "response", chatResponse)

	if chatResponse == nil {
		return false, fmt.Errorf("nil response from model")
	}

	if chatResponse.FinishReason == "" {
		return false, fmt.Errorf("response has no finish reason")
	}

	msg := chat.Message{
		Role:      chat.RoleAssistant,
		Content:   chatResponse.Content,
		Reasoning: chatResponse.Reasoning,
		ToolCalls: chatResponse.ToolCalls,
	}
	runResponse.Messages = append(runResponse.Messages, msg)
	chatResponse.Metrics.SessionID = runResponse.SessionID

	if runResponse.Metrics == nil {
		runResponse.Metrics = make(map[uuid.UUID][]metrics.Metrics)
	}
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
				return false, fmt.Errorf("tool call result not found for ID: %s", id)
			}
			runResponse.Messages = append(runResponse.Messages, chat.Message{
				Role:    chat.RoleTool,
				Content: result.Content,
				Name:    result.ToolCall.Name,
				IsErr:   result.IsErr,
			})
		}

		return true, nil
	} else {
		return false, nil
	}
}
