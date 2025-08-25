package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/yourlogarithm/l337/run"
)

type Tool struct {
	// Argument agnostic function wrapper over the user's implementaion.
	Callable ToolCallable
	// Parameters Schema
	Schema *jsonschema.Schema
	// SkillCard card
	SkillCard
}

type ToolCallable func(ctx context.Context, response *run.Response, rawArguments string) (string, error)

type ToolCallableTyped[T any] func(ctx context.Context, response *run.Response, args T) (string, error)

func wrapCallable[T any](fn ToolCallableTyped[T]) ToolCallable {
	return func(ctx context.Context, response *run.Response, rawArguments string) (string, error) {
		var args T
		if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
			return "", err
		}
		return fn(ctx, response, args)
	}
}

// Declare tool that does not require any arguments
func NewTool(name, description string, callable func(ctx context.Context) (string, error), options ...SkillCardOption) Tool {
	skill := SkillCard{
		Name:        name,
		Description: description,
	}
	for _, opt := range options {
		opt.Apply(&skill)
	}

	return Tool{
		Callable: func(ctx context.Context, response *run.Response, rawArguments string) (string, error) {
			return callable(ctx)
		},
		SkillCard: skill,
	}
}

// Declare a tool with required arguments
func NewToolWithArgs[T any](name, description string, callable ToolCallableTyped[T], options ...SkillCardOption) (Tool, error) {
	schema := jsonschema.Reflect(new(T))
	targetRef := strings.TrimPrefix(schema.Ref, "#/$defs/")
	if targetRef != "" {
		v, ok := schema.Definitions[targetRef]
		if !ok {
			return Tool{}, fmt.Errorf("definition %s not found", targetRef)
		}
		schema.Items = v.Items
		schema.Properties = v.Properties
		schema.Required = v.Required
		schema.Type = v.Type
		delete(schema.Definitions, targetRef)
	}

	skill := SkillCard{
		Name:        name,
		Description: description,
	}
	for _, opt := range options {
		opt.Apply(&skill)
	}

	return Tool{
		Callable:  wrapCallable(callable),
		SkillCard: skill,
		Schema:    schema,
	}, nil
}
