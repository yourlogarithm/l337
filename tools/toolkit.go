package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Helper map for managing tools.
type Toolkit map[string]Tool

func (t *Toolkit) Get(name string) (Tool, bool) {
	tool, exists := (*t)[name]
	return tool, exists
}

func (t *Toolkit) AddTool(tool Tool) {
	if *t == nil {
		*t = make(Toolkit)
	}
	(*t)[tool.Name] = tool
}

func (t *Toolkit) RegisterMCP(ctx context.Context, session *mcp.ClientSession) (int, error) {
	listToolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		return 0, err
	}
	for _, tool := range listToolsResult.Tools {
		t.AddTool(convertMCPTool(session, tool))
	}
	return len(listToolsResult.Tools), nil
}
