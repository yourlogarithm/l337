package anthropic

import (
	"context"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/yourlogarithm/l337/chat"
	internal_chat "github.com/yourlogarithm/l337/internal/chat"
	"github.com/yourlogarithm/l337/provider"
)

type anthropicProvider struct {
	model  anthropic.Model
	client anthropic.Client
}

func NewModel(name anthropic.Model, opts ...option.RequestOption) *provider.Model {
	return &provider.Model{
		Name:     string(name),
		Provider: "anthropic",
		Impl:     &anthropicProvider{model: name, client: anthropic.NewClient(opts...)},
	}
}

func (a *anthropicProvider) Chat(ctx context.Context, request *internal_chat.Request) (response internal_chat.Response, err error) {
	params := anthropic.MessageNewParams{
		Messages: make([]anthropic.MessageParam, 0, len(request.Messages)),
		Model:    a.model,
		Tools:    make([]anthropic.ToolUnionParam, 0, len(request.Tools)),
	}

	for _, msg := range request.Messages {
		var anthropicMsg anthropic.MessageParam
		switch msg.Role {
		case chat.RoleAssistant:
			anthropicMsg.Role = anthropic.MessageParamRoleAssistant
			anthropicMsg.Content = append(anthropicMsg.Content, anthropic.NewTextBlock(msg.Content))
		case chat.RoleUser:
			anthropicMsg.Role = anthropic.MessageParamRoleUser
			anthropicMsg.Content = append(anthropicMsg.Content, anthropic.NewTextBlock(msg.Content))
		default:
			return response, provider.NewUnknownRoleError(msg.Role.String())
		}
		params.Messages = append(params.Messages, anthropicMsg)
	}

	message, err := a.client.Messages.New(ctx, params)
	if err != nil {
		return response, err
	}

	return response, err
}
