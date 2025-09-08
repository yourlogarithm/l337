package ollama

import (
	"encoding/json"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func buildChatRequest(o *ollamaProvider, request *provider.Request, options *chat.Options, stream bool) (*api.ChatRequest, error) {
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: make([]api.Message, len(request.Messages)),
		Stream:   &stream,
		Tools:    make([]api.Tool, 0, len(request.Tools)),
		Options:  make(map[string]any),
	}

	if options.ReasoningEffort != nil {
		req.Think = &api.ThinkValue{Value: options.ReasoningEffort.AsAny()}
	}

	if options.ResponseFormat != nil {
		format, err := json.Marshal(options.ResponseFormat)
		if err != nil {
			return nil, err
		}
		req.Format = format
	}

	if options.KeepAlive != nil {
		req.KeepAlive = &api.Duration{Duration: *options.KeepAlive}
	}

	if options.Temperature != nil {
		req.Options["temperature"] = *options.Temperature
	}

	if options.Seed != nil {
		req.Options["seed"] = *options.Seed
	}

	if len(options.Stop) > 0 {
		req.Options["stop"] = options.Stop
	}

	if options.TopK != nil {
		req.Options["top_k"] = *options.TopK
	}

	if options.TopP != nil {
		req.Options["top_p"] = *options.TopP
	}

	for i, msg := range request.Messages {
		req.Messages[i].Role = msg.Role.String()
		if msg.Role == chat.RoleTool {
			req.Messages[i].ToolName = msg.Name
		}
		req.Messages[i].Content = msg.Content
		if len(msg.ToolCalls) > 0 {
			req.Messages[i].ToolCalls = make([]api.ToolCall, len(msg.ToolCalls))
			toolCalls := req.Messages[i].ToolCalls
			for j, toolCall := range msg.ToolCalls {
				id_int, err := strconv.Atoi(toolCall.ID)
				if err != nil {
					id_int = j
				}
				arguments := make(map[string]any)
				if err := json.Unmarshal([]byte(toolCall.Arguments), &arguments); err != nil {
					return nil, err
				}
				toolCalls[j] = api.ToolCall{
					Function: api.ToolCallFunction{
						Name:      toolCall.Name,
						Arguments: arguments,
						Index:     id_int,
					},
				}
			}
		}
	}

	for i := range request.Tools {
		ollamaTool := convertTool(&request.Tools[i])
		req.Tools = append(req.Tools, ollamaTool)
	}

	return req, nil
}
