package ollama

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/chat"
	internal_chat "github.com/yourlogarithm/l337/internal/chat"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/provider"
)

var logger = logging.SetupLogger("provider.ollama")

type ollamaProvider struct {
	model  string
	client *api.Client
}

func NewModel(name string, baseUrl string, http *http.Client) (*provider.Model, error) {
	baseUrlParsed, err := url.Parse(baseUrl)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(baseUrlParsed, http)
	return &provider.Model{
		Name:     name,
		Provider: "ollama",
		Impl:     &ollamaProvider{model: name, client: client},
	}, nil
}

func (o *ollamaProvider) Chat(ctx context.Context, request *internal_chat.Request, options *provider.ChatOptions) (response internal_chat.Response, err error) {
	stream := false
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: make([]api.Message, len(request.Messages)),
		Stream:   &stream,
		Tools:    make([]api.Tool, 0, len(request.Tools)),
		Options:  make(map[string]any),
	}

	if options.ReasoningEffort != nil {
		req.Think = &api.ThinkValue{Value: options.ReasoningEffort.AsAny()}
	}

	if options.ResponseFormat != nil {
		req.Format, err = json.Marshal(options.ResponseFormat)
		if err != nil {
			return response, err
		}
	}

	if options.KeepAlive != nil {
		req.KeepAlive = &api.Duration{Duration: *options.KeepAlive}
	}

	if options.Temperature != nil {
		req.Options["temperature"] = *options.Temperature
	}

	if options.Seed != nil {
		req.Options["seed"] = *options.Seed
	}

	if len(options.Stop) > 0 {
		req.Options["stop"] = options.Stop
	}

	if options.TopK != nil {
		req.Options["top_k"] = *options.TopK
	}

	if options.TopP != nil {
		req.Options["top_p"] = *options.TopP
	}

	for i, msg := range request.Messages {
		req.Messages[i].Role = msg.Role.String()
		if msg.Role == chat.RoleTool {
			req.Messages[i].ToolName = msg.Name
		}
		req.Messages[i].Content = msg.Content
		if len(msg.ToolCalls) > 0 {
			req.Messages[i].ToolCalls = make([]api.ToolCall, len(msg.ToolCalls))
			toolCalls := req.Messages[i].ToolCalls
			for j, toolCall := range msg.ToolCalls {
				id_int, err := strconv.Atoi(toolCall.ID)
				if err != nil {
					id_int = j
				}
				arguments := make(map[string]any)
				if err := json.Unmarshal([]byte(toolCall.Arguments), &arguments); err != nil {
					return response, err
				}
				toolCalls[j] = api.ToolCall{
					Function: api.ToolCallFunction{
						Name:      toolCall.Name,
						Arguments: arguments,
						Index:     id_int,
					},
				}
			}
		}
	}

	for i := range request.Tools {
		ollamaTool := convertTool(&request.Tools[i])
		req.Tools = append(req.Tools, ollamaTool)
	}

	callback := func(ollamaResp api.ChatResponse) error {
		logger.Debug("chat.response", "model", o.model, "response", ollamaResp)
		response.FinishReason = ollamaResp.DoneReason
		response.Content += ollamaResp.Message.Content

		metrics := convertMetrics(&ollamaResp.Metrics)
		response.Metrics.Add(&metrics)

		for _, toolCall := range ollamaResp.Message.ToolCalls {
			rawArguments, err := json.Marshal(toolCall.Function.Arguments)
			if err != nil {
				return err
			}

			response.ToolCalls = append(response.ToolCalls, chat.ToolCall{
				ID:        strconv.Itoa(toolCall.Function.Index),
				Arguments: string(rawArguments),
				Name:      toolCall.Function.Name,
			})
		}
		return nil
	}

	logger.Debug("chat.request", "model", o.model, "messages", request.Messages, "tools", request.Tools)

	if err = o.client.Chat(ctx, req, callback); err != nil {
		return internal_chat.Response{}, err
	}

	return response, nil
}
