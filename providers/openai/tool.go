package openai

import (
	"encoding/json"

	"github.com/openai/openai-go"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func convertTool(t *tools.Tool) (tool openai.ChatCompletionToolParam) {
	tool.Type = "function"

	tool.Function.Name = t.Name
	tool.Function.Description = openai.String(t.Description)

	if t.Parameters != nil {
		marshaled, _ := json.Marshal(t.Parameters)
		json.Unmarshal(marshaled, &tool.Function.Parameters)
	}

	return tool
}

func convertToolCall(toolCall *types.ToolCall) openai.ChatCompletionMessageToolCallParam {
	return openai.ChatCompletionMessageToolCallParam{
		ID: toolCall.ID,
		Function: openai.ChatCompletionMessageToolCallFunctionParam{
			Name:      toolCall.Name,
			Arguments: toolCall.Arguments,
		},
	}
}
