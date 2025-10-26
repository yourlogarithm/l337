package openai

import (
	"context"
	"errors"
	"time"

	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func (o *openAIProvider) ChatStreaming(
	ctx context.Context,
	messages []types.Message,
	tools []tools.Tool,
	options *types.Options,
) (streaming.ResponseChannel, error) {
	params, err := buildChatRequest(o.model, messages, tools, options)
	if err != nil {
		return nil, err
	}

	logger.Debug("chat.request.stream", "model", o.model, "messages", messages, "tools", tools)

	openaiStream := o.client.Chat.Completions.NewStreaming(ctx, params)
	channel := streaming.NewResponseChannel(options.StreamingBufferSize)

	go func() {
		defer channel.Close()

		start := time.Now()
		for openaiStream.Next() {
			chunk := openaiStream.Current()

			if len(chunk.Choices) == 0 && chunk.Usage.TotalTokens == 0 {
				channel.SendErr(errors.New("no choices in response"))
				return
			} else if len(chunk.Choices) == 0 {
				response := types.Response{
					Metrics: convertMetrics(&chunk.Usage, time.Since(start)),
				}
				channel.Send(&response)
				continue
			}

			choice := chunk.Choices[0]
			response := types.Response{
				ID:      chunk.ID,
				Created: time.Unix(chunk.Created, 0),
				Content: types.NewTextContent(choice.Delta.Content),
				Refusal: choice.Delta.Refusal,
				// Reasoning:    choice.Delta.ReasoningContent,
				FinishReason: choice.FinishReason,
			}
			for _, toolCall := range choice.Delta.ToolCalls {
				tc := types.ToolCall{
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
