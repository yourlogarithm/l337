package ollama

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func (o *ollamaProvider) Chat(ctx context.Context, request *provider.Request, options *chat.Options) (response provider.Response, err error) {
	req, err := buildChatRequest(o, request, options, false)
	if err != nil {
		return response, err
	}

	callback := func(ollamaResp api.ChatResponse) error {
		logger.Debug("chat.response", "model", o.model, "response", ollamaResp)

		response.Created = ollamaResp.CreatedAt
		response.FinishReason = ollamaResp.DoneReason
		response.Reasoning = ollamaResp.Message.Thinking
		response.Content = ollamaResp.Message.Content
		response.Metrics = convertMetrics(&ollamaResp.Metrics)

		for _, toolCall := range ollamaResp.Message.ToolCalls {
			rawArguments, err := json.Marshal(toolCall.Function.Arguments)
			if err != nil {
				return err
			}
			toolCall := chat.ToolCall{
				ID:        strconv.Itoa(toolCall.Function.Index),
				Arguments: string(rawArguments),
				Name:      toolCall.Function.Name,
			}
			response.ToolCalls = append(response.ToolCalls, toolCall)
		}

		return nil
	}

	logger.Debug("chat.request", "model", o.model, "messages", request.Messages, "tools", request.Tools)
	if err = o.client.Chat(ctx, req, callback); err != nil {
		logger.Error("chat.response", "model", o.model, "error", err)
		return provider.Response{}, err
	}

	return response, nil
}
