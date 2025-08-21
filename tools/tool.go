package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
)

type Tool struct {
	// Argument agnostic function wrapper over the user's implementaion.
	Callable ToolCallable

	Name        string
	Description string

	// Params schema
	Schema *jsonschema.Schema
}

type ToolCallable func(ctx context.Context, rawArguments string) (string, error)

type ToolCallableTyped[T any] func(ctx context.Context, args T) (string, error)

func wrapCallable[T any](fn ToolCallableTyped[T]) ToolCallable {
	return func(ctx context.Context, rawArguments string) (string, error) {
		var args T
		if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
			return "", err
		}
		return fn(ctx, args)
	}
}

func NewTool(name, description string, callable func(ctx context.Context) (string, error)) Tool {
	return Tool{
		Callable: func(ctx context.Context, rawArguments string) (string, error) {
			return callable(ctx)
		},
		Name:        name,
		Description: description,
	}
}

func NewToolWithArgs[T any](name, description string, callable ToolCallableTyped[T]) (Tool, error) {

	schema := jsonschema.Reflect(new(T))
	targetRef := strings.TrimPrefix(schema.Ref, "#/$defs/")
	v, ok := schema.Definitions[targetRef]
	if !ok {
		return Tool{}, fmt.Errorf("definition %s not found", targetRef)
	}
	schema.Items = v.Items
	schema.Properties = v.Properties
	schema.Required = v.Required
	schema.Type = v.Type
	delete(schema.Definitions, targetRef)

	return Tool{
		Callable: wrapCallable(callable),

		Name:        name,
		Description: description,
		Schema:      schema,
	}, nil
}
