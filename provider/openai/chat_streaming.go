package openai

import (
	"context"
	"errors"
	"time"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func (o *openAIProvider) ChatStreaming(
	ctx context.Context,
	request *provider.Request,
	options *chat.Options,
) (provider.ResponseChannel, error) {
	params, err := buildChatRequest(o.model, request, options)
	if err != nil {
		return nil, err
	}

	logger.Debug("chat.request.stream", "model", o.model, "messages", request.Messages, "tools", request.Tools)

	openaiStream := o.client.Chat.Completions.NewStreaming(ctx, params)
	channel := provider.NewResponseChannel(options.StreamingBufferSize)

	go func() {
		defer channel.Close()

		start := time.Now()
		for openaiStream.Next() {
			chunk := openaiStream.Current()

			if len(chunk.Choices) == 0 && chunk.Usage.TotalTokens == 0 {
				channel.SendErr(errors.New("no choices in response"))
				return
			} else if len(chunk.Choices) == 0 {
				response := provider.Response{
					Metrics: convertMetrics(&chunk.Usage, time.Since(start)),
				}
				channel.Send(&response)
				continue
			}

			choice := chunk.Choices[0]
			response := provider.Response{
				ID:      chunk.ID,
				Created: time.Unix(chunk.Created, 0),
				Content: choice.Delta.Content,
				Refusal: choice.Delta.Refusal,
				// Reasoning:    choice.Delta.ReasoningContent,
				FinishReason: choice.FinishReason,
			}
			for _, toolCall := range choice.Delta.ToolCalls {
				tc := chat.ToolCall{
					ID:        toolCall.ID,
					Name:      toolCall.Function.Name,
					Arguments: toolCall.Function.Arguments,
				}
				response.ToolCalls = append(response.ToolCalls, tc)
			}
			channel.Send(&response)
		}

		if err := openaiStream.Err(); err != nil {
			channel.SendErr(err)
		}
	}()

	return channel, nil
}
