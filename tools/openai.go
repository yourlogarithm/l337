package tools

import (
	"github.com/openai/openai-go"
)

func (t *Tool) ToOpenAITool() (tool openai.ChatCompletionToolParam) {
	tool.Type = "function"

	tool.Function.Name = t.Name
	tool.Function.Description = openai.String(t.Description)
	tool.Function.Parameters = make(map[string]any, len(t.Parameters))
	for name, schema := range t.Parameters {
		tool.Function.Parameters[name] = schema
	}

	return tool
}
