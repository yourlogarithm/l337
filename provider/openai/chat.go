package openai

import (
	"context"
	"time"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func (o *openAIProvider) Chat(ctx context.Context, request *provider.Request, options *chat.Options) (response provider.Response, err error) {
	params, err := buildChatRequest(o.model, request, options)
	if err != nil {
		return response, err
	}

	logger.Debug("chat.request", "model", o.model, "messages", request.Messages, "tools", request.Tools)
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

	response.Content = choice.Message.Content
	response.Refusal = choice.Message.Refusal
	response.ToolCalls = make([]chat.ToolCall, len(choice.Message.ToolCalls))
	response.FinishReason = choice.FinishReason
	response.Metrics = convertMetrics(&chatCompletion.Usage, totalDuration)

	for j, toolCall := range choice.Message.ToolCalls {
		response.ToolCalls[j] = chat.ToolCall{
			ID:        toolCall.ID,
			Arguments: toolCall.Function.Arguments,
			Name:      toolCall.Function.Name,
		}
	}

	return response, nil
}
