package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/yourlogarithm/l337/tools"
)

func convertTool(t *tools.Tool) (tool anthropic.ToolUnionParam) {
	inputSchema := anthropic.ToolInputSchemaParam{}
	if t.Schema != nil {
		inputSchema.Properties = t.Schema.Properties
		inputSchema.Required = t.Schema.Required
	}
	return anthropic.ToolUnionParamOfTool(inputSchema, t.Name)
}
