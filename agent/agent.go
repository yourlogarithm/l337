package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/provider"
	"github.com/yourlogarithm/golagno/retry"
	"github.com/yourlogarithm/golagno/run"
	"github.com/yourlogarithm/golagno/tools"
)

type Agent struct {
	Name        string
	Role        string
	Description string

	Instructions   string
	Goal           string
	ExpectedOutput string

	Model *provider.Model

	Tools tools.Toolkit

	RetryOptions *retry.Options
}

func (a *Agent) computeSystemMessage() chat.Message {
	var sb strings.Builder

	appendSystemString := func(s string, tag string) {
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteRune('\n')
			}
			if tag != "" {
				sb.WriteString("<" + tag + ">\n")
			}
			sb.WriteString(s)
			if tag != "" {
				sb.WriteRune('\n')
				sb.WriteString("</" + tag + ">")
			}
		}
	}

	appendSystemString(a.Description, "")
	appendSystemString(a.Goal, "goal")
	appendSystemString(a.Instructions, "instructions")
	appendSystemString(a.ExpectedOutput, "expected_output")

	return chat.Message{
		Role:    "system",
		Content: sb.String(),
	}
}

func (a *Agent) Run(ctx context.Context, messages []chat.Message) (runResponse run.Response, err error) {
	if a.RetryOptions == nil {
		a.RetryOptions = retry.Default()
	}

	// Generate system message
	runResponse.Messages = append(runResponse.Messages, a.computeSystemMessage())

	// Add user messages
	runResponse.Messages = append(runResponse.Messages, messages...)

	tools := make([]tools.Tool, 0, len(a.Tools))
	for _, tool := range a.Tools {
		tools = append(tools, tool)
	}

	for {
		var chatResponse chat.Response
		req := chat.Request{
			Messages: runResponse.Messages,
			Tools:    tools,
		}
		if err = a.RetryOptions.Execute(func() error {
			response, err := a.Model.Impl.Chat(ctx, &req)
			if err != nil {
				return err
			}
			chatResponse = response
			return nil
		}); err != nil {
			return runResponse, err
		}
		msg := chat.Message{
			Role:      "assistant",
			Content:   chatResponse.Content,
			ToolCalls: chatResponse.ToolCalls,
		}
		runResponse.Messages = append(runResponse.Messages, msg)
		if len(chatResponse.ToolCalls) > 0 {
			var wg sync.WaitGroup
			wg.Add(len(chatResponse.ToolCalls))

			type ToolCallResult struct {
				ToolCall *chat.ToolCall
				Content  string
			}

			results := make(map[string]ToolCallResult, len(chatResponse.ToolCalls))
			order := make([]string, 0, len(chatResponse.ToolCalls))

			for _, toolCall := range chatResponse.ToolCalls {
				order = append(order, toolCall.ID)
				go func(ctx context.Context, a *Agent, toolCall *chat.ToolCall) {
					defer wg.Done()
					result, err := a.Tools.Call(ctx, toolCall.Name, toolCall.Arguments)
					var content string
					if err != nil {
						content = "error: " + err.Error()
					} else {
						content = result
					}
					results[toolCall.ID] = ToolCallResult{
						ToolCall: toolCall,
						Content:  content,
					}
				}(ctx, a, &toolCall)
			}
			wg.Wait()

			for _, id := range order {
				result, exists := results[id]
				if !exists {
					return runResponse, fmt.Errorf("tool call result not found for ID: %s", id)
				}
				runResponse.Messages = append(runResponse.Messages, chat.Message{
					Role:    "tool",
					Content: result.Content,
					Name:    result.ToolCall.Name,
				})
			}
		} else if chatResponse.FinishReason == chat.FinishReasonStop {
			break
		}
	}

	return runResponse, nil
}
