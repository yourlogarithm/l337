package tools

import (
	"context"
	"fmt"
)

type Toolkit map[string]Tool

func (t *Toolkit) AddTool(tool Tool) {
	(*t)[tool.Name] = tool
}

func (t *Toolkit) Call(ctx context.Context, name string, args map[string]any) (string, error) {
	tool, exists := (*t)[name]
	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return tool.Callable(ctx, args)
}
