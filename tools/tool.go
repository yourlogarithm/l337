package tools

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/yourlogarithm/l337/types"
)

type Tool struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Parameters  *jsonschema.Schema `json:"parameters,omitempty"`

	// Argument agnostic function wrapper over the user's implementaion.
	callable ToolCallable `json:"-"`
}

// Declare tool that does not require any arguments
func New(name, description string, callable func(ctx context.Context) ([]types.Content, error)) Tool {
	return Tool{
		Name:        name,
		Description: description,
		callable: func(ctx context.Context, response *types.Run, rawArguments string) ([]types.Content, error) {
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

func (t Tool) Call(ctx context.Context, run *types.Run, rawArguments string) ([]types.Content, error) {
	return t.callable(ctx, run, rawArguments)
}
