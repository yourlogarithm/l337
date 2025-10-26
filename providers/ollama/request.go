package ollama

import (
	"encoding/json"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/providers"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func buildChatRequest(o *ollamaProvider, messages []types.Message, tools []tools.Tool, options *types.Options, stream bool) (*api.ChatRequest, error) {
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: make([]api.Message, len(messages)),
		Stream:   &stream,
		Tools:    make([]api.Tool, 0, len(tools)),
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

	for i, msg := range messages {
		req.Messages[i].Role = msg.Role.String()
		if msg.Role == types.RoleTool {
			req.Messages[i].ToolName = msg.Name
		}
		req.Messages[i].Content = msg.Content.Text
		if msg.Content.Image != nil {
			if msg.Content.Image.Url != "" {
				return nil, providers.ErrParams{Param: "Message.Content.Image", Msg: "ollama api does not support image url message content, use bytes buffer instead"}
			}
			req.Messages[i].Images = append(req.Messages[i].Images, api.ImageData(msg.Content.Image.ImageData.Base64))
		}
		if msg.Content.Audio != nil {
			return nil, providers.ErrParams{Param: "Message.Content.Audio", Msg: "ollama api does not support audio message content"}
		}
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

	for i := range tools {
		ollamaTool := convertTool(&tools[i])
		req.Tools = append(req.Tools, ollamaTool)
	}

	return req, nil
}
