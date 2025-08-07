package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/logging"
	"github.com/yourlogarithm/l337/provider"
)

var logger = logging.SetupLogger("provider.openai")

type OpenAI struct {
	model  string
	client openai.Client
}

func NewOpenAI(name string, opts ...option.RequestOption) *provider.Model {
	return &provider.Model{
		Name:     name,
		Provider: "openai",
		Impl:     &OpenAI{model: name, client: openai.NewClient(opts...)},
	}
}

func (o *OpenAI) Chat(ctx context.Context, request *chat.Request) (response chat.Response, err error) {
	params := openai.ChatCompletionNewParams{
		Messages: make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Messages)),
		Model:    o.model,
		Tools:    make([]openai.ChatCompletionToolParam, 0, len(request.Tools)),
	}

	for _, msg := range request.Messages {
		var openaiMsg openai.ChatCompletionMessageParamUnion
		switch msg.Role {
		case chat.RoleDeveloper:
			openaiMsg = openai.DeveloperMessage(msg.Content)
			openaiMsg.OfDeveloper.Name = openai.String(msg.Name)
		case chat.RoleSystem:
			openaiMsg = openai.SystemMessage(msg.Content)
			openaiMsg.OfSystem.Name = openai.String(msg.Name)
		case chat.RoleUser:
			openaiMsg = openai.UserMessage(msg.Content)
			openaiMsg.OfUser.Name = openai.String(msg.Name)
		case chat.RoleAssistant:
			openaiMsg = openai.AssistantMessage(msg.Content)
			openaiMsg.OfAssistant.Name = openai.String(msg.Name)
		case chat.RoleTool:
			openaiMsg = openai.ToolMessage(msg.Content, msg.Name)
		case chat.RoleFunction:
			openaiMsg = openai.ChatCompletionMessageParamOfFunction(msg.Content, msg.Name)
		default:
			return response, provider.NewUnknownRoleError(msg.Role.String())
		}
		params.Messages = append(params.Messages, openaiMsg)
	}

	for i := range request.Tools {
		params.Tools = append(params.Tools, convertTool(&request.Tools[i]))
	}

	logger.Debug("chat.request", "model", o.model, "messages", request.Messages, "tools", request.Tools)
	chatCompletion, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return response, err
	}
	logger.Debug("chat.response", "model", o.model, "response", chatCompletion)

	response.ID = chatCompletion.ID
	response.Created = chatCompletion.Created
	choice := chatCompletion.Choices[0]

	response.Content = choice.Message.Content
	response.Refusal = choice.Message.Refusal
	response.ToolCalls = make([]chat.ToolCall, len(choice.Message.ToolCalls))
	response.FinishReason = chat.FinishReason(choice.FinishReason)

	for j, toolCall := range choice.Message.ToolCalls {
		arguments := make(map[string]any)
		if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &arguments); err != nil {
			return response, fmt.Errorf("OpenAI Chat: failed to unmarshal tool call arguments: %w", err)
		}

		response.ToolCalls[j] = chat.ToolCall{
			ID:        toolCall.ID,
			Arguments: arguments,
			Name:      toolCall.Function.Name,
		}
	}

	return response, nil
}
