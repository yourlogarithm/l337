package openai

import (
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"github.com/openai/openai-go/shared"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func buildChatRequest(model string, request *provider.Request, options *chat.Options) (openai.ChatCompletionNewParams, error) {
	params := openai.ChatCompletionNewParams{
		Messages:         make([]openai.ChatCompletionMessageParamUnion, 0, len(request.Messages)),
		Model:            model,
		Tools:            make([]openai.ChatCompletionToolParam, 0, len(request.Tools)),
		Logprobs:         param.NewOpt(options.Logprobs),
		PresencePenalty:  param.NewOpt(options.PresencePenalty),
		PromptCacheKey:   param.NewOpt(options.PromptCacheKey),
		SafetyIdentifier: param.NewOpt(options.SafetyIdentifier),
		User:             param.NewOpt(options.User),
		LogitBias:        options.LogitBias,
		ServiceTier:      openai.ChatCompletionNewParamsServiceTier(options.ServiceTier),
		Stop:             openai.ChatCompletionNewParamsStopUnion{OfStringArray: options.Stop},
	}

	if len(options.Modalities) > 0 {
		modalities := make([]string, len(options.Modalities))
		for i, modality := range options.Modalities {
			modalities[i] = string(modality)
		}
		params.Modalities = modalities
	}

	if options.Audio.Voice != "" || options.Audio.Format != "" {
		params.Audio = openai.ChatCompletionAudioParam{
			Voice:  openai.ChatCompletionAudioParamVoice(options.Audio.Voice),
			Format: openai.ChatCompletionAudioParamFormat(options.Audio.Format),
		}
	}

	if options.MaxCompletionTokens > 0 {
		params.MaxCompletionTokens = param.NewOpt(int64(options.MaxCompletionTokens))
	}

	if options.ReasoningEffort != nil {
		if level, ok := options.ReasoningEffort.AsLevel(); ok {
			params.ReasoningEffort = shared.ReasoningEffort(level)
		} else {
			return params, provider.ErrParams{Param: "ReasoningEffort", Msg: fmt.Sprintf("invalid reasoning effort: %v", options.ReasoningEffort)}
		}
	}
	if options.MaxTokens > 0 {
		params.MaxTokens = param.NewOpt(int64(options.MaxTokens))
	}
	if options.FrequencyPenalty != nil {
		params.FrequencyPenalty = param.NewOpt(*options.FrequencyPenalty)
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
	if options.IncludeStreamMetrics {
		params.StreamOptions.IncludeUsage = param.NewOpt(true)
	}

	for _, msg := range request.Messages {
		var openaiMsg openai.ChatCompletionMessageParamUnion
		var nameParam *param.Opt[string]
		switch msg.Role {
		case chat.RoleDeveloper:
			developerMsg := openai.ChatCompletionDeveloperMessageParam{}
			if msg.Content.Text != "" {
				developerMsg.Content.OfString = param.NewOpt(msg.Content.Text)
			} else {
				return params, provider.ErrParams{Param: "Message.Content", Msg: "developer messages must have non-empty text content"}
			}
			openaiMsg = openai.ChatCompletionMessageParamUnion{OfDeveloper: &developerMsg}
			nameParam = &openaiMsg.OfDeveloper.Name
		case chat.RoleSystem:
			systemMsg := openai.ChatCompletionSystemMessageParam{}
			if msg.Content.Text != "" {
				systemMsg.Content.OfString = param.NewOpt(msg.Content.Text)
			} else {
				return params, provider.ErrParams{Param: "Message.Content", Msg: "system messages must have non-empty text content"}
			}
			openaiMsg = openai.ChatCompletionMessageParamUnion{OfSystem: &systemMsg}
			nameParam = &openaiMsg.OfSystem.Name
		case chat.RoleUser:
			userMsg := openai.ChatCompletionUserMessageParam{}
			if msg.Content.Text != "" {
				userMsg.Content.OfString = param.NewOpt(msg.Content.Text)
			} else if msg.Content.Audio != nil {
				userMsg.Content.OfArrayOfContentParts = append(userMsg.Content.OfArrayOfContentParts, openai.ChatCompletionContentPartUnionParam{
					OfInputAudio: &openai.ChatCompletionContentPartInputAudioParam{
						InputAudio: openai.ChatCompletionContentPartInputAudioInputAudioParam{
							Data:   msg.Content.Audio.Base64,
							Format: string(msg.Content.Audio.Format),
						},
					},
				})
			} else if msg.Content.Image != nil {
				var url strings.Builder
				if msg.Content.Image.ImageData != nil {
					url.WriteString("data:image/")
					url.WriteString(string(msg.Content.Image.ImageData.Format))
					url.WriteString(";base64,")
					url.WriteString(msg.Content.Image.ImageData.Base64)
				} else if msg.Content.Image.Url != "" {
					url.WriteString(msg.Content.Image.Url)
				}
				userMsg.Content.OfArrayOfContentParts = append(userMsg.Content.OfArrayOfContentParts, openai.ChatCompletionContentPartUnionParam{
					OfImageURL: &openai.ChatCompletionContentPartImageParam{
						ImageURL: openai.ChatCompletionContentPartImageImageURLParam{
							URL:    url.String(),
							Detail: string(msg.Content.Image.Detail),
						},
					},
				})
			} else {
				return params, provider.ErrParams{Param: "Message.Content", Msg: "user messages must have non-empty text, audio, or image content"}
			}
			openaiMsg = openai.ChatCompletionMessageParamUnion{OfUser: &userMsg}
			nameParam = &openaiMsg.OfUser.Name
		case chat.RoleAssistant:
			assistantMsg := openai.ChatCompletionAssistantMessageParam{}
			if msg.Content.Text != "" {
				assistantMsg.Content.OfString = param.NewOpt(msg.Content.Text)
			} else {
				return params, provider.ErrParams{Param: "Message.Content", Msg: "assistant messages must have non-empty text content"}
			}
			if len(msg.ToolCalls) > 0 {
				assistantMsg.ToolCalls = make([]openai.ChatCompletionMessageToolCallParam, len(msg.ToolCalls))
				for i := range msg.ToolCalls {
					assistantMsg.ToolCalls[i] = convertToolCall(&msg.ToolCalls[i])
				}
			}
			openaiMsg = openai.ChatCompletionMessageParamUnion{OfAssistant: &assistantMsg}
			nameParam = &openaiMsg.OfAssistant.Name
		case chat.RoleTool:
			toolMsg := openai.ChatCompletionToolMessageParam{
				ToolCallID: msg.Name,
			}
			if msg.Content.Text != "" {
				toolMsg.Content.OfString = param.NewOpt(msg.Content.Text)
			} else {
				return params, provider.ErrParams{Param: "Message.Content", Msg: "tool messages must have non-empty text content"}
			}
			openaiMsg = openai.ChatCompletionMessageParamUnion{OfTool: &toolMsg}
		default:
			return params, provider.ErrUnknownRole{Role: msg.Role.String()}
		}
		if msg.Name != "" && nameParam != nil {
			*nameParam = param.NewOpt(msg.Name)
		}
		params.Messages = append(params.Messages, openaiMsg)
	}

	for i := range request.Tools {
		params.Tools = append(params.Tools, convertTool(&request.Tools[i]))
	}

	return params, nil
}
