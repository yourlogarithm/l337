package ollama

import (
	"context"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func (o *ollamaProvider) ChatStreaming(
	ctx context.Context,
	request *provider.Request,
	options *chat.Options,
) (provider.ResponseChannel, error) {
	req, err := buildChatRequest(o, request, options, true)
	if err != nil {
		return nil, err
	}

	logger.Debug("chat.request.stream", "model", o.model, "messages", request.Messages, "tools", request.Tools)
	channel := provider.NewResponseChannel(options.StreamingBufferSize)

	go func() {
		defer channel.Close()

		ollamaCallback := func(ollamaResp api.ChatResponse) error {
			logger.Debug("chat.response.stream", "model", o.model, "response", ollamaResp)

			response := provider.Response{
				Created:      ollamaResp.CreatedAt,
				FinishReason: ollamaResp.DoneReason,
				Reasoning:    ollamaResp.Message.Thinking,
				Content:      chat.NewTextContent(ollamaResp.Message.Content),
				Metrics:      convertMetrics(&ollamaResp.Metrics),
			}

			for _, toolCall := range ollamaResp.Message.ToolCalls {
				convertedToolCall, err := convertToolCall(&toolCall)
				if err != nil {
					return err
				}
				response.ToolCalls = append(response.ToolCalls, convertedToolCall)
			}

			channel.Send(&response)

			return nil
		}

		err := o.client.Chat(ctx, req, ollamaCallback)
		if err != nil {
			channel.SendErr(err)
		}
	}()

	return channel, nil
}
