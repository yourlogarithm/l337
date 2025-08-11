package anthropic

import (
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/yourlogarithm/l337/tools"
)

func convertTool(t *tools.Tool) (tool anthropic.ToolUnionParam) {
	inputSchema := anthropic.ToolInputSchemaParam{
		Required:   t.Required,
		Properties: t.Parameters,
	}
	return anthropic.ToolUnionParamOfTool(inputSchema, t.Name)
}
