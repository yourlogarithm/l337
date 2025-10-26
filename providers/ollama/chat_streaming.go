package ollama

import (
	"context"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func (o *ollamaProvider) ChatStreaming(
	ctx context.Context,
	messages []types.Message,
	tools []tools.Tool,
	options *types.Options,
) (streaming.ResponseChannel, error) {
	req, err := buildChatRequest(o, messages, tools, options, true)
	if err != nil {
		return nil, err
	}

	logger.Debug("chat.request.stream", "model", o.model, "messages", messages, "tools", tools)
	channel := streaming.NewResponseChannel(options.StreamingBufferSize)

	go func() {
		defer channel.Close()

		ollamaCallback := func(ollamaResp api.ChatResponse) error {
			logger.Debug("chat.response.stream", "model", o.model, "response", ollamaResp)

			response := types.Response{
				Created:      ollamaResp.CreatedAt,
				FinishReason: ollamaResp.DoneReason,
				Reasoning:    ollamaResp.Message.Thinking,
				Content:      types.NewTextContent(ollamaResp.Message.Content),
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
