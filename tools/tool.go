package tools

import (
	"context"
	"reflect"
)

type Tool struct {
	Callable    ToolCallable
	Name        string
	Description string
	Parameters  []Parameter
}

type ToolCallable func(ctx context.Context, toolParams Params) (string, error)

func NewTool(name, description string, callable ToolCallable) Tool {
	return Tool{
		Name:        name,
		Description: description,
		Callable:    callable,
	}
}

func (t *Tool) AddParameter(p Parameter) {
	t.Parameters = append(t.Parameters, p)
}

func AddParameterFromType[T any](tool *Tool, name string) error {
	var zero T
	element, err := getElement(reflect.TypeOf(zero))
	if err != nil {
		return err
	}
	tool.Parameters = append(tool.Parameters, Parameter{
		Name:    name,
		Element: element,
	})
	return nil
}
