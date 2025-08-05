package tools_test

import (
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go/shared"
	"github.com/stretchr/testify/assert"
	"github.com/yourlogarithm/golagno/tools"
)

func TestToOpenAITool_BasicDataTypes(t *testing.T) {
	tool := tools.Tool{
		Name:        "simple_tool",
		Description: "A tool with simple parameters",
		Parameters: map[string]jsonschema.Schema{
			"param1": *jsonschema.Reflect(""),
			"param2": *jsonschema.Reflect(0),
		},
	}

	openaiTool := tool.ToOpenAITool()
	assert.Equal(t, "simple_tool", openaiTool.Function.Name)
	assert.Equal(t, "A tool with simple parameters", openaiTool.Function.Description.Value)

	expected := shared.FunctionParameters{
		"param1": *jsonschema.Reflect(""),
		"param2": *jsonschema.Reflect(0),
	}

	assert.Equal(t, openaiTool.Function.Parameters, expected)
}
