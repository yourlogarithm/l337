package provider

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/golagno/chat"
)

type Ollama struct {
	model  string
	client *api.Client
}

func NewOllama(name string, base *url.URL, http *http.Client) *Model {
	client := api.NewClient(base, http)
	return &Model{
		Name:     name,
		Provider: "ollama",
		Impl:     &Ollama{model: name, client: client},
	}
}

func (o *Ollama) Chat(ctx context.Context, messages []chat.Message) (response chat.Response, err error) {
	stream := false
	req := &api.ChatRequest{
		Model:    o.model,
		Messages: make([]api.Message, len(messages)),
		Stream:   &stream,
	}
	for i, msg := range messages {
		req.Messages[i].Role = msg.Role
		req.Messages[i].Content = msg.Content
	}

	callback := func(ollamaResp api.ChatResponse) error {
		slog.Debug("ollama.Chat.Response", "model", o.model, "ollamaResp", ollamaResp)
		if response.Choices == nil {
			response.Choices = make([]chat.Choice, 1)
		}
		choice := &response.Choices[0]
		choice.Content += ollamaResp.Message.Content
		choice.FinishReason = chat.FinishReason(ollamaResp.DoneReason)
		return nil
	}

	slog.Debug("ollama.Chat", "model", o.model, "messages", messages)

	if err = o.client.Chat(ctx, req, callback); err != nil {
		return chat.Response{}, err
	}

	return response, nil
}
