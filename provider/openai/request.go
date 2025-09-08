package openai

import (
	"fmt"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func buildChatRequest(model string, request *provider.Request, options *chat.Options) (openai.ChatCompletionNewParams, error) {
	params := openai.ChatCompletionNewParams{
		Messages:            make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Messages)),
		Model:               model,
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
			return params, ErrParams{Param: "ReasoningEffort", Msg: fmt.Sprintf("invalid reasoning effort: %v", options.ReasoningEffort)}
		}
	}
	if options.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(int64(options.MaxTokens))
	}
	if options.FrequencyPenalty != nil {
		params.FrequencyPenalty = param.NewOpt(*options.FrequencyPenalty)
	}
	// if options.N != nil {
	// 	params.N = param.NewOpt(int64(*options.N))
	// }
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
	if options.IncludeStreamMetrics {
		params.StreamOptions.IncludeUsage = param.NewOpt(true)
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
			return params, provider.ErrUnknownRole{Role: msg.Role.String()}
		}
		params.Messages = append(params.Messages, openaiMsg)
	}

	for i := range request.Tools {
		params.Tools = append(params.Tools, convertTool(&request.Tools[i]))
	}

	return params, nil
}
