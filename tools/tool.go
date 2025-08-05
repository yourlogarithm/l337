package tools

import (
	"context"

	"github.com/invopop/jsonschema"
)

type Tool struct {
	Callable    ToolCallable
	Name        string
	Description string
	Parameters  map[string]jsonschema.Schema
}

type ToolCallable func(ctx context.Context, toolParams Params) (string, error)

func NewTool(name, description string, callable ToolCallable) Tool {
	return Tool{
		Name:        name,
		Description: description,
		Callable:    callable,
		Parameters:  make(map[string]jsonschema.Schema),
	}
}

func AddParameterFromType[T any](tool *Tool, name string, description string) {
	var zero T
	schema := jsonschema.Reflect(zero)
	if description != "" {
		schema.Description = description
	}
	tool.Parameters[name] = *schema
}
