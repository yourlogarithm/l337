package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SayHiParams struct {
	Name string `json:"name"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, args SayHiParams) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Hi " + args.Name},
		},
	}, nil, nil
}

func main() {
	ctx := context.Background()
	_, serverTransport := mcp.NewInMemoryTransports()

	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v0.0.1"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		log.Fatal(err)
	}

	serverSession.Wait()
}
