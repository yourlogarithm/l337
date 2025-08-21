package openai

import (
	"context"
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/yourlogarithm/l337/chat"
	internal_chat "github.com/yourlogarithm/l337/internal/chat"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/provider"
)

var logger = logging.SetupLogger("provider.openai")

type openAIProvider struct {
	model  string
	client openai.Client
}

func NewModel(name string, opts ...option.RequestOption) *provider.Model {
	return &provider.Model{
		Name:     name,
		Provider: "openai",
		Impl:     &openAIProvider{model: name, client: openai.NewClient(opts...)},
	}
}

func (o *openAIProvider) Chat(ctx context.Context, request *internal_chat.Request, options *provider.ChatOptions) (response internal_chat.Response, err error) {
	params := openai.ChatCompletionNewParams{
		Messages:            make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Messages)),
		Model:               o.model,
		Tools:               make([]openai.ChatCompletionToolParam, 0, len(request.Tools)),
		Logprobs:            param.NewOpt(options.Logprobs),
		MaxCompletionTokens: param.NewOpt(int64(options.MaxCompletionTokens)),
		PresencePenalty:     param.NewOpt(options.PresencePenalty),
		PromptCacheKey:      param.NewOpt(options.PromptCacheKey),
		SafetyIdentifier:    param.NewOpt(options.SafetyIdentifier),
		User:                param.NewOpt(options.User),
		LogitBias:           options.LogitBias,
		ServiceTier:         openai.ChatCompletionNewParamsServiceTier(options.ServiceTier),
		Stop:                openai.ChatCompletionNewParamsStopUnion{OfStringArray: options.Stop},
	}

	if options.ReasoningEffort != nil {
		if level, ok := options.ReasoningEffort.AsLevel(); ok {
			params.ReasoningEffort = shared.ReasoningEffort(level)
		} else {
			return response, fmt.Errorf("invalid reasoning effort: %v", options.ReasoningEffort)
		}
	}

	if options.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(int64(options.MaxTokens))
	}

	if options.FrequencyPenalty != nil {
		params.FrequencyPenalty = param.NewOpt(*options.FrequencyPenalty)
	}

	if options.N != nil {
		params.N = param.NewOpt(int64(*options.N))
	}

	if options.Seed != nil {
		params.Seed = param.NewOpt(int64(*options.Seed))
	}

	if options.Temperature != nil {
		params.Temperature = param.NewOpt(*options.Temperature)
	}

	if options.TopLogprobs != nil {
		params.TopLogprobs = param.NewOpt(int64(*options.TopLogprobs))
	}

	if options.TopP != nil {
		params.TopP = param.NewOpt(*options.TopP)
	}

	if options.ParallelToolCalls != nil {
		params.ParallelToolCalls = param.NewOpt(*options.ParallelToolCalls)
	}

	if options.ResponseFormat != nil {
		params.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{JSONSchema: shared.ResponseFormatJSONSchemaJSONSchemaParam{Schema: options.ResponseFormat}}}
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
	response.FinishReason = choice.FinishReason

	for j, toolCall := range choice.Message.ToolCalls {
		response.ToolCalls[j] = chat.ToolCall{
			ID:        toolCall.ID,
			Arguments: toolCall.Function.Arguments,
			Name:      toolCall.Function.Name,
		}
	}

	return response, nil
}
