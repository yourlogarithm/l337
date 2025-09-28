package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yourlogarithm/l337/chat"
)

func convertMCPTool(session *mcp.ClientSession, mcpTool *mcp.Tool) (Tool, error) {
	schema, err := json.Marshal(mcpTool.InputSchema)
	if err != nil {
		return Tool{}, ErrToolCreation{Message: "failed to marshal `mcpTool.InputSchema`"}
	}

	t := Tool{
		Name:        mcpTool.Name,
		Description: mcpTool.Description,
		callable: func(ctx context.Context, runResponse *chat.RunResponse, rawArguments string) (string, error) {
			args := map[string]any{}
			if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
				return "", fmt.Errorf("invalid tool arguments: %w", err)
			}
			params := &mcp.CallToolParams{
				Name:      mcpTool.Name,
				Arguments: args,
			}
			callResult, err := session.CallTool(ctx, params)
			if err != nil {
				return "", err
			} else if callResult.IsError {
				return "", fmt.Errorf("tool call error: %s", callResult.Content)
			}
			var sb strings.Builder
			for i, content := range callResult.Content {
				if textContent, ok := content.(*mcp.TextContent); ok {
					sb.WriteString(textContent.Text)
					if i > 0 {
						sb.WriteByte('\n')
					}
				} else {
					return "", fmt.Errorf("unsupported content type: %T", content)
				}
			}
			return sb.String(), nil
		},
	}

	if err := json.Unmarshal(schema, &t.Parameters); err != nil {
		return Tool{}, ErrToolCreation{Message: "faild to unmarshal tool parameters"}
	}

	return t, nil
}
