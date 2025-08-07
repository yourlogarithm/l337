package ollama_test

// func TestToOllamaTool_BasicDataTypes(t *testing.T) {
// 	tool := tools.Tool{
// 		Name:        "simple_tool",
// 		Description: "A tool with simple parameters",
// 		Parameters: map[string]jsonschema.Schema{
// 			"param1": *jsonschema.Reflect(""),
// 			"param2": *jsonschema.Reflect(0),
// 		},
// 	}

// 	ollamaTool := ollama.convertTool(&tool)
// 	assert.Equal(t, "simple_tool", ollamaTool.Function.Name)
// 	assert.Equal(t, "A tool with simple parameters", ollamaTool.Function.Description)
// 	assert.Equal(t, "object", ollamaTool.Function.Parameters.Type)

// 	assert.Len(t, ollamaTool.Function.Parameters.Properties, 2)
// 	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param1")
// 	assert.Contains(t, ollamaTool.Function.Parameters.Properties, "param2")

// 	assert.Equal(t, api.PropertyType{"string"}, ollamaTool.Function.Parameters.Properties["param1"].Type)
// 	assert.Empty(t, ollamaTool.Function.Parameters.Properties["param1"].Items)
// 	assert.Equal(t, api.PropertyType{"integer"}, ollamaTool.Function.Parameters.Properties["param2"].Type)
// 	assert.Empty(t, ollamaTool.Function.Parameters.Properties["param2"].Items)
// }
