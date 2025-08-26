package google

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourlogarithm/l337/chat"
	internal_chat "github.com/yourlogarithm/l337/internal/chat"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/provider"
	"google.golang.org/genai"
)

var logger = logging.SetupLogger("provider.google")

type googleProvider struct {
	model  string
	client *genai.Client
}

func NewModel(ctx context.Context, name string, config *genai.ClientConfig) (*provider.Model, error) {
	client, err := genai.NewClient(ctx, config)
	if err != nil {
		return nil, err
	}
	return &provider.Model{
		Name:     name,
		Provider: "google",
		Impl:     &googleProvider{model: name, client: client},
	}, nil
}

func convertFloatPointer(fp64 *float64) *float32 {
	if fp64 == nil {
		return nil
	}
	fp32 := float32(*fp64)
	return &fp32
}

func (g *googleProvider) Chat(ctx context.Context, request *internal_chat.Request, options *provider.ChatOptions) (response internal_chat.Response, err error) {

	config := &genai.GenerateContentConfig{
		TopK:             options.TopK,
		MaxOutputTokens:  int32(options.MaxCompletionTokens),
		StopSequences:    options.Stop,
		ResponseLogprobs: options.Logprobs,
		Logprobs:         options.TopLogprobs,
		ResponseMIMEType: options.ResponseMIMEType,
		CachedContent:    options.PromptCacheKey,
	}

	if options.Temperature != nil {
		config.Temperature = convertFloatPointer(options.Temperature)
	}
	if options.TopP != nil {
		config.TopP = convertFloatPointer(options.TopP)
	}
	// if options.N != nil {
	// 	config.CandidateCount = int32(*options.N)
	// }
	if options.PresencePenalty != 0 {
		config.PresencePenalty = convertFloatPointer(&options.PresencePenalty)
	}
	if options.FrequencyPenalty != nil {
		config.FrequencyPenalty = convertFloatPointer(options.FrequencyPenalty)
	}
	if options.Seed != nil {
		converted := int32(*options.Seed)
		config.Seed = &converted
	}
	if options.ResponseFormat != nil {
		converted, err := convertSchema(options.ResponseFormat)
		if err != nil {
			return response, err
		}
		config.ResponseSchema = &converted
		if config.ResponseMIMEType == "" {
			config.ResponseMIMEType = "application/json"
		} else if config.ResponseMIMEType != "application/json" {
			return response, fmt.Errorf("unsupported response MIME type: %s", config.ResponseMIMEType)
		}
	}
	for i := range request.Tools {
		converted, err := convertTool(&request.Tools[i])
		if err != nil {
			return response, err
		}
		config.Tools = append(config.Tools, &converted)
	}
	if options.Thinking > 0 {
		if config.ThinkingConfig == nil {
			config.ThinkingConfig = &genai.ThinkingConfig{}
		}
		converted := int32(options.Thinking)
		config.ThinkingConfig.ThinkingBudget = &converted
	}
	if options.ReasoningEffort != nil {
		if config.ThinkingConfig == nil {
			config.ThinkingConfig = &genai.ThinkingConfig{}
		}
		think, ok := options.ReasoningEffort.AsBool()
		if ok {
			config.ThinkingConfig.IncludeThoughts = think
		} else {
			return response, fmt.Errorf("reasoning effort must be a boolean for Google provider")
		}
	}

	var contents = make([]*genai.Content, 0, len(request.Messages))
	for _, msg := range request.Messages {
		switch msg.Role {
		case chat.RoleSystem:
			config.SystemInstruction = &genai.Content{
				Parts: []*genai.Part{{Text: msg.Content}},
			}
		case chat.RoleUser:
			contents = append(contents, &genai.Content{
				Parts: []*genai.Part{{Text: msg.Content}},
				Role:  msg.Role.String(),
			})
		case chat.RoleAssistant, chat.RoleModel:
			content := genai.Content{
				Role: chat.RoleModel.String(),
			}
			if msg.Content != "" && (len(msg.ToolCalls) == 0 || len(msg.ToolCalls) == 1) {
				content.Parts = []*genai.Part{{Text: msg.Content}}
				for i := range msg.ToolCalls {
					converted, err := convertToolCall(&msg.ToolCalls[i])
					if err != nil {
						return response, err
					}
					content.Parts[0].FunctionCall = &converted
				}
			} else {
				if msg.Content != "" {
					content.Parts = append(content.Parts, &genai.Part{Text: msg.Content})
				}
				for i := range msg.ToolCalls {
					converted, err := convertToolCall(&msg.ToolCalls[i])
					if err != nil {
						return response, err
					}
					content.Parts = append(content.Parts, &genai.Part{FunctionCall: &converted})
				}
			}
			contents = append(contents, &content)
		case chat.RoleTool:
			functionResponse := genai.FunctionResponse{
				Name:     msg.Name,
				Response: make(map[string]any, 1),
			}
			if msg.IsErr {
				functionResponse.Response["error"] = msg.Content
			} else {
				functionResponse.Response["output"] = msg.Content
			}
			content := genai.Content{
				Parts: []*genai.Part{
					{
						FunctionResponse: &functionResponse,
					},
				},
				Role: chat.RoleUser.String(),
			}
			contents = append(contents, &content)
		default:
			return response, fmt.Errorf("unsupported message role: %s", msg.Role)
		}
	}

	logger.Debug("chat.request", "model", g.model, "messages", request.Messages, "tools", request.Tools, "options", options)
	start := time.Now()
	result, err := g.client.Models.GenerateContent(
		ctx,
		g.model,
		contents,
		config,
	)
	totalDuration := time.Since(start)
	if err != nil {
		return response, err
	}
	logger.Debug("chat.response", "model", g.model, "response", result)

	response.Metrics = convertMetrics(result.UsageMetadata, totalDuration)

	if len(result.Candidates) == 0 {
		return response, fmt.Errorf("no candidates returned")
	}

	candidate := result.Candidates[0]

	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			var buffer *string
			if part.Thought {
				buffer = &response.Reasoning
			} else {
				buffer = &response.Content
			}
			if *buffer != "" {
				*buffer += "\n"
			}
			*buffer += part.Text
		}
		if part.FunctionCall != nil {
			marshaled, err := json.Marshal(part.FunctionCall.Args)
			if err != nil {
				return response, err
			}
			toolCall := chat.ToolCall{
				ID:        part.FunctionCall.ID,
				Arguments: string(marshaled),
				Name:      part.FunctionCall.Name,
			}
			response.ToolCalls = append(response.ToolCalls, toolCall)
		}
	}

	response.FinishReason = string(candidate.FinishReason)

	return response, nil
}
