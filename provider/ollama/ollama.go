package ollama

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/logging"
	"github.com/yourlogarithm/golagno/provider"
)

var logger = logging.SetupLogger("provider.ollama")

type Ollama struct {
	model  string
	client *api.Client
}

func NewOllama(name string, base *url.URL, http *http.Client) *provider.Model {
	client := api.NewClient(base, http)
	return &provider.Model{
		Name:     name,
		Provider: "ollama",
		Impl:     &Ollama{model: name, client: client},
	}
}

func (o *Ollama) Chat(ctx context.Context, request *chat.Request) (response chat.Response, err error) {
	stream := false
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: make([]api.Message, len(request.Messages)),
		Stream:   &stream,
		Tools:    make([]api.Tool, 0, len(request.Tools)),
	}

	for i, msg := range request.Messages {
		req.Messages[i].Role = msg.Role
		req.Messages[i].Content = msg.Content
	}

	for _, tool := range request.Tools {
		ollamaTool := tool.ToOllamaTool()
		req.Tools = append(req.Tools, ollamaTool)
	}

	callback := func(ollamaResp api.ChatResponse) error {
		logger.Debug("chat.response", "model", o.model, "response", ollamaResp)
		response.FinishReason = chat.FinishReason(ollamaResp.DoneReason)
		response.Content += ollamaResp.Message.Content
		for _, toolCall := range ollamaResp.Message.ToolCalls {
			response.ToolCalls = append(response.ToolCalls, chat.ToolCall{
				ID:        strconv.Itoa(toolCall.Function.Index),
				Arguments: toolCall.Function.Arguments,
				Name:      toolCall.Function.Name,
			})
		}
		return nil
	}

	logger.Debug("chat.request", "model", o.model, "messages", request.Messages, "tools", request.Tools)

	if err = o.client.Chat(ctx, req, callback); err != nil {
		return chat.Response{}, err
	}

	return response, nil
}
