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
	Required    []string

	// ModifiesRunResponse indicates if this tool modifies the RunResponse during `.Run()`
	// If so *RunResponse will be injected into `toolParams` for use.
	ModifiesRunResponse bool
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

func AddParameterFromType[T any](tool *Tool, name string, description string, required bool) {
	var zero T
	schema := jsonschema.Reflect(zero)
	if description != "" {
		schema.Description = description
	}
	tool.Parameters[name] = *schema
	if required {
		tool.Required = append(tool.Required, name)
	}
}
