package openai

import (
	"context"
	"time"

	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func (o *openAIProvider) Chat(ctx context.Context, messages []types.Message, tools []tools.Tool, options *types.Options) (response types.Response, err error) {
	params, err := buildChatRequest(o.model, messages, tools, options)
	if err != nil {
		return response, err
	}

	logger.Debug("chat.request", "model", o.model, "messages", messages, "tools", tools)
	start := time.Now()
	chatCompletion, err := o.client.Chat.Completions.New(ctx, params)
	totalDuration := time.Since(start)
	if err != nil {
		return response, err
	}
	logger.Debug("chat.response", "model", o.model, "response", chatCompletion)

	response.ID = chatCompletion.ID
	response.Created = time.Unix(chatCompletion.Created, 0)
	choice := chatCompletion.Choices[0]

	response.Content = types.NewTextContent(choice.Message.Content)
	if choice.Message.Audio.Data != "" {
		response.Content.Audio = &types.Audio{
			Base64: choice.Message.Audio.Data,
		}
	}
	response.Refusal = choice.Message.Refusal
	response.ToolCalls = make([]types.ToolCall, len(choice.Message.ToolCalls))
	response.FinishReason = choice.FinishReason
	response.Metrics = convertMetrics(&chatCompletion.Usage, totalDuration)

	for j, toolCall := range choice.Message.ToolCalls {
		response.ToolCalls[j] = types.ToolCall{
			ID:        toolCall.ID,
			Arguments: toolCall.Function.Arguments,
			Name:      toolCall.Function.Name,
		}
	}

	return response, nil
}
