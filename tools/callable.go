package tools

import (
	"context"
	"encoding/json"

	"github.com/yourlogarithm/l337/types"
)

type ToolCallable func(ctx context.Context, run *types.Run, rawArguments string) ([]types.Content, error)

type ToolCallableTyped[T any] func(context.Context, *types.Run, T) ([]types.Content, error)

func wrapCallable[T any](fn ToolCallableTyped[T]) ToolCallable {
	return func(ctx context.Context, response *types.Run, rawArguments string) ([]types.Content, error) {
		var args T
		if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
			return nil, err
		}
		return fn(ctx, response, args)
	}
}
