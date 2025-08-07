package openai

import (
	"maps"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
	"github.com/yourlogarithm/l337/tools"
)

func convertTool(t *tools.Tool) (tool openai.ChatCompletionToolParam) {
	tool.Type = "function"

	tool.Function.Name = t.Name
	tool.Function.Description = openai.String(t.Description)
	tool.Function.Parameters = make(map[string]any, len(t.Parameters))

	properties := make(map[string]jsonschema.Schema, len(t.Parameters))
	maps.Copy(properties, t.Parameters)
	tool.Function.Parameters["properties"] = properties
	if len(t.Required) > 0 {
		tool.Function.Parameters["required"] = t.Required
	}

	return tool
}
