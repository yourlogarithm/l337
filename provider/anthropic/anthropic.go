package anthropic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
			for _, toolCall := range msg.ToolCalls {
				anthropicMsg.Content = append(anthropicMsg.Content, anthropic.NewToolUseBlock(toolCall.ID, toolCall.Arguments, toolCall.Name))
			}
			params.Messages = append(params.Messages, anthropicMsg)
		case chat.RoleUser:
			anthropicMsg.Role = anthropic.MessageParamRoleUser
			params.Messages = append(params.Messages, anthropicMsg)
		case chat.RoleSystem:
			params.System = append(params.System, anthropic.TextBlockParam{Text: msg.Content})
		case chat.RoleTool:
			anthropicMsg.Role = anthropic.MessageParamRoleUser
			anthropicMsg.Content = append(anthropicMsg.Content, anthropic.NewToolResultBlock(msg.Name, msg.Content, msg.IsErr))
			params.Messages = append(params.Messages, anthropicMsg)
		default:
			return response, provider.NewUnknownRoleError(msg.Role.String())
		}
	}

	for i := range request.Tools {
		toolParam := convertTool(&request.Tools[i])
		params.Tools = append(params.Tools, toolParam)
	}

	message, err := a.client.Messages.New(ctx, params)
	if err != nil {
		return response, err
	}

	response.ID = message.ID
	response.Created = time.Now().Unix()
	response.FinishReason = string(message.StopReason)

	for _, contentBlock := range message.Content {
		switch contentBlock.Type {
		case "text":
			response.Content += contentBlock.Text
		case "tool_use":
			toolCall := chat.ToolCall{
				ID:   contentBlock.ID,
				Name: contentBlock.Name,
			}
			if err := json.Unmarshal(contentBlock.Input, &toolCall.Arguments); err != nil {
				return response, err
			}
			response.ToolCalls = append(response.ToolCalls, toolCall)
		default:
			return response, fmt.Errorf("unsupported content block type: %s", contentBlock.Type)
		}
	}

	return response, err
}
