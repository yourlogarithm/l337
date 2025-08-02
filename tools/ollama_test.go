package tools_test

import (
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/ollama/ollama/api"
	"github.com/stretchr/testify/assert"
	"github.com/yourlogarithm/golagno/tools"
	"github.com/yourlogarithm/golagno/tools/test"
)

func TestToOllamaTool_SimpleParameters(t *testing.T) {
	tool := tools.Tool{
		Name:        "simple_tool",
		Description: "A tool with simple parameters",
		Parameters: []tools.Parameter{
			{Name: "param1", Element: tools.Element{Type: "string"}},
			{Name: "param2", Element: tools.Element{Type: "number"}},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)
	assert.Equal(t, "simple_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with simple parameters", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	assert.Len(t, ollamaTool.Function.Parameters.Properties, 2)
	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param1")
	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param2")

	assert.Equal(t, api.PropertyType{"string"}, ollamaTool.Function.Parameters.Properties["param1"].Type)
	assert.Empty(t, ollamaTool.Function.Parameters.Properties["param1"].Items)
	assert.Equal(t, api.PropertyType{"number"}, ollamaTool.Function.Parameters.Properties["param2"].Type)
	assert.Empty(t, ollamaTool.Function.Parameters.Properties["param2"].Items)
}

func TestToOllamaTool_StructParameter(t *testing.T) {
	tool := tools.Tool{
		Name:        "struct_tool",
		Description: "A tool with a struct parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type: "object",
					Schema: &jsonschema.Schema{
						Title: "StructSchema",
					},
				},
			},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)
	assert.Equal(t, "struct_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with a struct parameter", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	assert.Len(t, ollamaTool.Function.Parameters.Properties, 1)
	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param1")
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Properties["param1"].Type)
}

func TestToOllamaTool_StructParameterWithSchema(t *testing.T) {
	expectedSchema := jsonschema.Reflect(test.Foo{})

	tool := tools.Tool{
		Name:        "struct_tool",
		Description: "A tool with a struct parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type:   "Foo",
					Schema: expectedSchema,
				},
			},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)

	// Validate tool metadata
	assert.Equal(t, "struct_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with a struct parameter", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	// Validate properties
	assert.Len(t, ollamaTool.Function.Parameters.Properties, 1)
	param1 := ollamaTool.Function.Parameters.Properties["param1"]
	assert.Equal(t, "Foo", param1.Type.String())
	assert.Empty(t, param1.Items)

	assert.Equal(t, []jsonschema.Schema{*expectedSchema}, ollamaTool.Function.Parameters.Defs)
}

func TestToOllamaTool_SliceParameter(t *testing.T) {
	tool := tools.Tool{
		Name:        "slice_tool",
		Description: "A tool with a slice parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type: "array",
					Nested: &tools.Element{
						Type: "string",
					},
				},
			},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)
	assert.Equal(t, "slice_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with a slice parameter", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	assert.Len(t, ollamaTool.Function.Parameters.Properties, 1)
	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param1")
	assert.Equal(t, api.PropertyType{"array"}, ollamaTool.Function.Parameters.Properties["param1"].Type)
	assert.Equal(t, "string", ollamaTool.Function.Parameters.Properties["param1"].Items)
}

func TestToOllamaTool_SliceParameterWithNestedSchema(t *testing.T) {
	expectedSchema := jsonschema.Reflect(test.Foo{})

	tool := tools.Tool{
		Name:        "slice_tool",
		Description: "A tool with a slice parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type: "array",
					Nested: &tools.Element{
						Type:   "Foo",
						Schema: expectedSchema,
					},
				},
			},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)

	// Validate tool metadata
	assert.Equal(t, "slice_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with a slice parameter", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	// Validate properties
	assert.Len(t, ollamaTool.Function.Parameters.Properties, 1)
	param1 := ollamaTool.Function.Parameters.Properties["param1"]
	assert.Equal(t, "array", param1.Type.String())
	assert.Equal(t, "Foo", param1.Items)

	assert.Equal(t, []jsonschema.Schema{*expectedSchema}, ollamaTool.Function.Parameters.Defs)
}

func TestToOllamaTool_MapParameter(t *testing.T) {
	tool := tools.Tool{
		Name:        "map_tool",
		Description: "A tool with a map parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type: "object",
					Nested: &tools.Element{
						Type: "string",
					},
				},
			},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)
	assert.Equal(t, "map_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with a map parameter", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	assert.Len(t, ollamaTool.Function.Parameters.Properties, 1)
	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param1")
	assert.Equal(t, api.PropertyType{"object"}, ollamaTool.Function.Parameters.Properties["param1"].Type)
	assert.Equal(t, "string", ollamaTool.Function.Parameters.Properties["param1"].Items)
}

func TestToOllamaTool_MapParameterWithNestedSchema(t *testing.T) {
	expectedSchema := jsonschema.Reflect(test.Foo{})
	tool := tools.Tool{
		Name:        "map_tool",
		Description: "A tool with a map parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type: "object",
					Nested: &tools.Element{
						Type:   "Foo",
						Schema: expectedSchema,
					},
				},
			},
		},
	}

	ollamaTool, err := tool.ToOllamaTool()
	assert.NoError(t, err)

	// Validate tool metadata
	assert.Equal(t, "map_tool", ollamaTool.Function.Name)
	assert.Equal(t, "A tool with a map parameter", ollamaTool.Function.Description)
	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

	// Validate properties
	assert.Len(t, ollamaTool.Function.Parameters.Properties, 1)
	param1 := ollamaTool.Function.Parameters.Properties["param1"]
	assert.Equal(t, "object", param1.Type.String())
	assert.Equal(t, "Foo", param1.Items)

	assert.Equal(t, []jsonschema.Schema{*expectedSchema}, ollamaTool.Function.Parameters.Defs)
}

func TestToOllamaTool_NestedArrayParameter(t *testing.T) {
	tool := tools.Tool{
		Name:        "nested_array_tool",
		Description: "A tool with a nested array parameter",
		Parameters: []tools.Parameter{
			{
				Name: "param1",
				Element: tools.Element{
					Type: "array",
					Nested: &tools.Element{
						Type: "array",
						Nested: &tools.Element{
							Type: "string",
						},
					},
				},
			},
		},
	}

	_, err := tool.ToOllamaTool()
	assert.Error(t, err)
	assert.EqualError(t, err, "ToOllamaTool: Ollama api does not support double nested 2D arrays or maps")
}
