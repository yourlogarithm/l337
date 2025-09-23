package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/yourlogarithm/l337/chat"
)

type Tool struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  *jsonschema.Schema `json:"parameters,omitempty"`

	// Argument agnostic function wrapper over the user's implementaion.
	callable ToolCallable `json:"-"`
}

func (t Tool) Call(ctx context.Context, runResponse *chat.RunResponse, rawArguments string) (string, error) {
	return t.callable(ctx, runResponse, rawArguments)
}

type ToolCallable func(ctx context.Context, runResponse *chat.RunResponse, rawArguments string) (string, error)

type ToolCallableTyped[T any] func(context.Context, *chat.RunResponse, T) (string, error)

func wrapCallable[T any](fn ToolCallableTyped[T]) ToolCallable {
	return func(ctx context.Context, response *chat.RunResponse, rawArguments string) (string, error) {
		var args T
		if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
			return "", err
		}
		return fn(ctx, response, args)
	}
}

// Declare tool that does not require any arguments
func NewTool(name, description string, callable func(ctx context.Context) (string, error)) Tool {
	return Tool{
		Name:        name,
		Description: description,
		callable: func(ctx context.Context, response *chat.RunResponse, rawArguments string) (string, error) {
			return callable(ctx)
		},
	}
}

// Declare a tool with required arguments
func NewToolWithArgs[T any](name, description string, callable ToolCallableTyped[T]) (Tool, error) {
	schema := jsonschema.Reflect(new(T))
	targetRef := strings.TrimPrefix(schema.Ref, "#/$defs/")
	if targetRef != "" {
		v, ok := schema.Definitions[targetRef]
		if !ok {
			return Tool{}, ErrToolCreation{Message: fmt.Sprintf("failed to find definition for %s", targetRef)}
		}
		schema.Items = v.Items
		schema.Properties = v.Properties
		schema.Required = v.Required
		schema.Type = v.Type
		delete(schema.Definitions, targetRef)
	}

	return Tool{
		Name:        name,
		Description: description,
		Parameters:  schema,
		callable:    wrapCallable(callable),
	}, nil
}
