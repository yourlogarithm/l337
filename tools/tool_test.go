package tools_test

import (
	"testing"

	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/yourlogarithm/golagno/tools"
	"github.com/yourlogarithm/golagno/tools/test"
)

func testing_function_name() {}

func TestCreateToolFromFunc_Name(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(testing_function_name)
	assert.Equal(t, "testing_function_name", tool.Name)
}

// single line description
func description_single_line() {}

func TestCreateToolFromFunc_Description_SingleLine(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(description_single_line)
	assert.Equal(t, "single line description", tool.Description)
}

// multi-line 0 description
// multi-line 1 description
func description_multi_line() {}

func TestCreateToolFromFunc_Description_MultiLine(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(description_multi_line)
	assert.Equal(t, "multi-line 0 description\nmulti-line 1 description", tool.Description)
}

/*
block-line 0 description
block-line 1 description
*/
func description_block() {}

func TestCreateToolFromFunc_Description_Block(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(description_block)
	assert.Equal(t, tool.Description, "block-line 0 description\nblock-line 1 description")
}

func single_base_parameter(x string) {}

func TestCreateToolFromFunc_SingleBaseParameter(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(single_base_parameter)
	assert.Len(t, tool.Parameters, 1)
	parameter := tool.Parameters[0]
	assert.Equal(t, "x", parameter.Name)
	assert.Equal(t, tools.Element{Type: "string"}, parameter.Element)
}

func multiple_base_parameters(x string, y int64, z bool) {}

func TestCreateToolFromFunc_MultipleBaseParameters(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(multiple_base_parameters)
	assert.Len(t, tool.Parameters, 3)
	assert.Equal(t, "x", tool.Parameters[0].Name)
	assert.Equal(t, tools.Element{Type: "string"}, tool.Parameters[0].Element)
	assert.Equal(t, "y", tool.Parameters[1].Name)
	assert.Equal(t, tools.Element{Type: "number"}, tool.Parameters[1].Element)
	assert.Equal(t, "z", tool.Parameters[2].Name)
	assert.Equal(t, tools.Element{Type: "bool"}, tool.Parameters[2].Element)
}

func single_struct_parameter(foo test.Foo) {}

func TestCreateToolFromFunc_SingleStructParameter(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(single_struct_parameter)
	assert.Len(t, tool.Parameters, 1)
	parameter := tool.Parameters[0]
	assert.Equal(t, "foo", parameter.Name)
	assert.Equal(t, "Foo", parameter.Element.Type)
	assert.Nil(t, parameter.Element.Nested)

	expectedSchema := jsonschema.Reflect(test.Foo{})
	assert.Equal(t, expectedSchema, parameter.Element.Schema)
}

func multiple_struct_parameters(foo test.Foo, bar test.Bar) {}

func TestCreateToolFromFunc_MultipleStructParameters(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(multiple_struct_parameters)
	assert.Len(t, tool.Parameters, 2)

	fooParam := tool.Parameters[0]
	assert.Equal(t, "foo", fooParam.Name)
	assert.Equal(t, "Foo", fooParam.Element.Type)
	assert.Nil(t, fooParam.Element.Nested)
	expectedFooSchema := jsonschema.Reflect(test.Foo{})
	assert.Equal(t, expectedFooSchema, fooParam.Element.Schema)

	barParam := tool.Parameters[1]
	assert.Equal(t, "bar", barParam.Name)
	assert.Equal(t, "Bar", barParam.Element.Type)
	assert.Nil(t, barParam.Element.Nested)
	expectedBarSchema := jsonschema.Reflect(test.Bar{})
	assert.Equal(t, expectedBarSchema, barParam.Element.Schema)
}

func slice_parameter(fooSlice []test.Foo, barSlice []test.Bar) {}

func TestCreateToolFromFunc_SliceParameter(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(slice_parameter)
	assert.Len(t, tool.Parameters, 2)

	fooParam := tool.Parameters[0]
	assert.Equal(t, "fooSlice", fooParam.Name)
	assert.Equal(t, "array", fooParam.Element.Type)
	assert.NotNil(t, fooParam.Element.Nested)
	assert.Equal(t, "Foo", fooParam.Element.Nested.Type)
	expectedFooSchema := jsonschema.Reflect(test.Foo{})
	assert.Equal(t, expectedFooSchema, fooParam.Element.Nested.Schema)

	barParam := tool.Parameters[1]
	assert.Equal(t, "barSlice", barParam.Name)
	assert.Equal(t, "array", barParam.Element.Type)
	assert.NotNil(t, barParam.Element.Nested)
	assert.Equal(t, "Bar", barParam.Element.Nested.Type)
	expectedBarSchema := jsonschema.Reflect(test.Bar{})
	assert.Equal(t, expectedBarSchema, barParam.Element.Nested.Schema)
}

func map_parameter(fooMap map[string]test.Foo, barMap map[string]test.Bar) {}

func TestCreateToolFromFunc_MapParameter(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(map_parameter)
	assert.Len(t, tool.Parameters, 2)

	fooParam := tool.Parameters[0]
	assert.Equal(t, "fooMap", fooParam.Name)
	assert.Equal(t, "object", fooParam.Element.Type)
	assert.NotNil(t, fooParam.Element.Nested)
	assert.Equal(t, "Foo", fooParam.Element.Nested.Type)
	expectedFooSchema := jsonschema.Reflect(test.Foo{})
	assert.Equal(t, expectedFooSchema, fooParam.Element.Nested.Schema)

	barParam := tool.Parameters[1]
	assert.Equal(t, "barMap", barParam.Name)
	assert.Equal(t, "object", barParam.Element.Type)
	assert.NotNil(t, barParam.Element.Nested)
	assert.Equal(t, "Bar", barParam.Element.Nested.Type)
	expectedBarSchema := jsonschema.Reflect(test.Bar{})
	assert.Equal(t, expectedBarSchema, barParam.Element.Nested.Schema)
}

func mixed_parameters(foo test.Foo, bar test.Bar, baz []string, qux map[string]int) {}

func TestCreateToolFromFunc_MixedParameters(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(mixed_parameters)
	assert.Len(t, tool.Parameters, 4)

	fooParam := tool.Parameters[0]
	assert.Equal(t, "foo", fooParam.Name)
	assert.Equal(t, "Foo", fooParam.Element.Type)
	expectedFooSchema := jsonschema.Reflect(test.Foo{})
	assert.Equal(t, expectedFooSchema, fooParam.Element.Schema)

	barParam := tool.Parameters[1]
	assert.Equal(t, "bar", barParam.Name)
	assert.Equal(t, "Bar", barParam.Element.Type)
	expectedBarSchema := jsonschema.Reflect(test.Bar{})
	assert.Equal(t, expectedBarSchema, barParam.Element.Schema)

	bazParam := tool.Parameters[2]
	assert.Equal(t, "baz", bazParam.Name)
	assert.Equal(t, "array", bazParam.Element.Type)
	assert.NotNil(t, bazParam.Element.Nested)
	assert.Equal(t, "string", bazParam.Element.Nested.Type)

	quxParam := tool.Parameters[3]
	assert.Equal(t, "qux", quxParam.Name)
	assert.Equal(t, "object", quxParam.Element.Type)
	assert.NotNil(t, quxParam.Element.Nested)
	assert.Equal(t, "number", quxParam.Element.Nested.Type)
}

func invalid_map_parameter(fooMap map[test.Bar]test.Foo) {}

func TestCreateToolFromFunc_InvalidMapParameter(t *testing.T) {
	_, err := tools.CreateToolFromFunc(invalid_map_parameter)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "map keys must be strings")
}

func nested_array(arr [][]string) {}

func TestCreateToolFromFunc_NestedArrayParameter(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(nested_array)
	assert.Len(t, tool.Parameters, 1)

	param := tool.Parameters[0]
	assert.Equal(t, "arr", param.Name)
	assert.Equal(t, "array", param.Element.Type)
	assert.NotNil(t, param.Element.Nested)
	assert.Equal(t, "array", param.Element.Nested.Type)
	assert.NotNil(t, param.Element.Nested.Nested)
	assert.Equal(t, "string", param.Element.Nested.Nested.Type)
}

func nested_map(nestedMap map[string]map[string]string) {}

func TestCreateToolFromFunc_NestedMapParameter(t *testing.T) {
	tool, _ := tools.CreateToolFromFunc(nested_map)
	assert.Len(t, tool.Parameters, 1)

	param := tool.Parameters[0]
	assert.Equal(t, "nestedMap", param.Name)
	assert.Equal(t, "object", param.Element.Type)
	assert.NotNil(t, param.Element.Nested)
	assert.Equal(t, "object", param.Element.Nested.Type)
	assert.NotNil(t, param.Element.Nested.Nested)
	assert.Equal(t, "string", param.Element.Nested.Nested.Type)
}
