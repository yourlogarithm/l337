package tools

import (
	"context"

	"github.com/invopop/jsonschema"
)

type Tool struct {
	// Function that implements the tool's functionality.
	Callable    ToolCallable
	Name        string
	Description string
	// Map of parameter names to their JSON schema definitions.
	Parameters map[string]jsonschema.Schema
	// List of mandatory parameters.
	Required []string
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

// Add a parameter by name, description and required status.
//
// The JSON schema for the parameter is inferred from the type `T`.
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
