package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/xeipuuv/gojsonschema"
	"github.com/yourlogarithm/l337/types"
)

func (t *Toolkit) RegisterMCP(ctx context.Context, session *mcp.ClientSession) (int, error) {
	listToolsResult, err := session.ListTools(ctx, nil)
	if err != nil {
		return 0, err
	}

	for _, tool := range listToolsResult.Tools {
		tt, err := convertMCPTool(session, tool)
		if err != nil {
			return 0, err
		}
		t.AddTool(tt)
	}

	return len(listToolsResult.Tools), nil
}

func convertMCPTool(session *mcp.ClientSession, mcpTool *mcp.Tool) (Tool, error) {
	inputSchema, err := json.Marshal(mcpTool.InputSchema)
	if err != nil {
		return Tool{}, ErrToolCreation{Message: "failed to marshal `mcpTool.InputSchema`"}
	}

	var outputLoader *gojsonschema.JSONLoader
	if mcpTool.OutputSchema != nil {
		outputSchemaBuff, err := json.Marshal(mcpTool.OutputSchema)
		if err != nil {
			return Tool{}, ErrToolCreation{Message: "failed to marshal `mcpTool.OutputSchema`"}
		}
		loader := gojsonschema.NewBytesLoader(outputSchemaBuff)
		outputLoader = &loader
	}

	t := Tool{
		Name:        mcpTool.Name,
		Description: mcpTool.Description,
		callable: func(ctx context.Context, run *types.Run, rawArguments string) ([]types.Content, error) {
			args := map[string]any{}
			if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
				return nil, err
			}
			params := &mcp.CallToolParams{
				Name:      mcpTool.Name,
				Arguments: args,
			}
			callResult, err := session.CallTool(ctx, params)
			if err != nil {
				return nil, err
			} else if callResult.IsError {
				return nil, fmt.Errorf("tool call error: %s", callResult.Content)
			}
			var sb strings.Builder
			for i, content := range callResult.Content {
				if textContent, ok := content.(*mcp.TextContent); ok {
					sb.WriteString(textContent.Text)
					if i > 0 {
						sb.WriteByte('\n')
					}
				} else {
					return nil, fmt.Errorf("unsupported content type: %T", content)
				}
			}

			// https://modelcontextprotocol.io/specification/2025-06-18/server/tools#:~:text=Clients%20SHOULD%20validate%20structured%20results%20against%20this%20schema.

			out := sb.String()

			if outputLoader != nil {
				result, err := gojsonschema.Validate(*outputLoader, gojsonschema.NewStringLoader(out))
				if err != nil {
					return nil, err
				}
				if !result.Valid() {
					return nil, fmt.Errorf("output does not conform to schema: %v", result.Errors())
				}
			}

			return types.NewTextContent(out).AsSlice(), nil
		},
	}

	if err := json.Unmarshal(inputSchema, &t.Parameters); err != nil {
		return Tool{}, ErrToolCreation{Message: "failed to unmarshal tool parameters"}
	}

	return t, nil
}
