package main

import (
	"context"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SayHiParams struct {
	Name string `json:"name"`
}

type SayHiResult struct {
	Greeting string `json:"greeting"`
}

func sayHi(ctx context.Context, req *mcp.CallToolRequest, args SayHiParams) (*mcp.CallToolResult, SayHiResult, error) {
	return nil, SayHiResult{Greeting: "Hi " + args.Name}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v0.0.1"}, nil)

	tool := mcp.Tool{
		Name:        "greet",
		Description: "say hi",
	}

	mcp.AddTool(server, &tool, sayHi)

	handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server { return server }, nil)

	log.Print("Starting server at :8080")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
