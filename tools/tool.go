package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
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

func (t Tool) Call(ctx context.Context, runResponse *chat.RunResponse, rawArguments string) ([]chat.Content, error) {
	return t.callable(ctx, runResponse, rawArguments)
}

type ToolCallable func(ctx context.Context, runResponse *chat.RunResponse, rawArguments string) ([]chat.Content, error)

type ToolCallableTyped[T any] func(context.Context, *chat.RunResponse, T) ([]chat.Content, error)

func wrapCallable[T any](fn ToolCallableTyped[T]) ToolCallable {
	return func(ctx context.Context, response *chat.RunResponse, rawArguments string) ([]chat.Content, error) {
		var args T
		if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
			return nil, err
		}
		return fn(ctx, response, args)
	}
}

// Declare tool that does not require any arguments
func New(name, description string, callable func(ctx context.Context) ([]chat.Content, error)) Tool {
	return Tool{
		Name:        name,
		Description: description,
		callable: func(ctx context.Context, response *chat.RunResponse, rawArguments string) ([]chat.Content, error) {
			return callable(ctx)
		},
	}
}

// Declare a tool with required arguments
// Arguments `T` must be a struct type, otherwise an error is returned.
func NewWithArgs[T any](name, description string, callable ToolCallableTyped[T]) (Tool, error) {
	var t T
	if reflect.TypeOf(t).Kind() != reflect.Struct {
		return Tool{}, ErrToolCreation{Message: "`callable ToolCallableTyped[T]` expects T to be a struct type"}
	}

	schema := jsonschema.Reflect(t)
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
	schema.Ref = ""
	schema.Version = ""

	return Tool{
		Name:        name,
		Description: description,
		Parameters:  schema,
		callable:    wrapCallable(callable),
	}, nil
}
