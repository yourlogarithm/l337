package provider

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/golagno/chat"
)

type OpenAI struct {
	model  string
	client openai.Client
}

func NewOpenAI(name string, opts ...option.RequestOption) *Model {
	return &Model{
		Name:     name,
		Provider: "openai",
		Impl:     &OpenAI{model: name, client: openai.NewClient(opts...)},
	}
}

func (o *OpenAI) Chat(ctx context.Context, request *chat.Request) (chat.Response, error) {
	response := chat.Response{}
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, len(request.Messages))

	for _, msg := range request.Messages {
		var openaiMsg openai.ChatCompletionMessageParamUnion
		switch msg.Role {
		case "developer":
			openaiMsg = openai.DeveloperMessage(msg.Content)
			openaiMsg.OfDeveloper.Name = openai.String(msg.Name)
		case "system":
			openaiMsg = openai.SystemMessage(msg.Content)
			openaiMsg.OfSystem.Name = openai.String(msg.Name)
		case "user":
			openaiMsg = openai.UserMessage(msg.Content)
			openaiMsg.OfUser.Name = openai.String(msg.Name)
		case "assistant":
			openaiMsg = openai.AssistantMessage(msg.Content)
			openaiMsg.OfAssistant.Name = openai.String(msg.Name)
		case "tool":
			// openaiMsg = openai.ToolMessage(msg.Content, msg.ToolCallID)
			return response, fmt.Errorf("tool message not implemented yet")
		case "function":
			openaiMsg = openai.ChatCompletionMessageParamOfFunction(msg.Content, msg.Name)
		default:
			return response, NewUnknownRoleError(msg.Role)
		}
	}

	chatCompletion, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openaiMessages,
		Model:    o.model,
	})
	if err != nil {
		return response, err
	}

	response.ID = chatCompletion.ID
	response.Created = chatCompletion.Created
	choice := chatCompletion.Choices[0]

	response.Content = choice.Message.Content
	response.Refusal = choice.Message.Refusal
	response.ToolCalls = make([]chat.ToolCall, len(choice.Message.ToolCalls))
	response.FinishReason = chat.FinishReason(choice.FinishReason)

	for j, toolCall := range choice.Message.ToolCalls {
		if len(toolCall.Function.Arguments) > 0 {
			return response, fmt.Errorf("OpenAI Chat: tool call arguments are not supported yet")
		}
		response.ToolCalls[j] = chat.ToolCall{
			ID: toolCall.ID,
			// Arguments: toolCall.Function.Arguments,
			Name: toolCall.Function.Name,
		}
	}

	return response, nil
}
