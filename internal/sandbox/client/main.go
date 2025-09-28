package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
)

func main() {
	llmHTTPClient, logger := newLoggingHTTPClient()

	model := openai.NewModel(
		"gpt-5-nano",
		// option.WithBaseURL(os.Getenv("BASE_URL")),
		option.WithAPIKey(os.Getenv("API_KEY")),
		option.WithHTTPClient(llmHTTPClient),
	)

	// model, _ := ollama.NewModel(
	// 	"qwen3:8b",
	// 	"http://localhost:11434",
	// 	llmHTTPClient,
	// )

	options := chat.Options{
		// ReasoningEffort: chat.NewReasoningEffortBool(true),
	}

	mcpClient := mcp.NewClient(&mcp.Implementation{}, nil)
	transport := mcp.SSEClientTransport{
		Endpoint:   "http://localhost:8080/",
		HTTPClient: http.DefaultClient,
	}
	session, err := mcpClient.Connect(context.Background(), &transport, nil)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	a, err := agent.New(model, agent.WithChatOptions(options), agent.WithMCP(context.Background(), session))
	if err != nil {
		panic(err)
	}

	runResponse, stream, err := a.RunStreamingWithParams(
		context.Background(),
		0,
		chat.WithMessage(chat.RoleUser, "Call greet concurrently with names Alice, Bob, and Charlie."),
	)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}

	fmt.Println("Consuming stream...")
	content := false
	reasoning := false
	toolCalls := false
	for {
		chunk, err := stream.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error occurred: %v\n", err)
			break
		}
		if len(chunk.Reasoning) > 0 {
			if !reasoning {
				fmt.Printf("\n\n--- Reasoning ---\n")
				reasoning = true
				content = false
				toolCalls = false
			}
			fmt.Print(chunk.Reasoning)
		}
		if len(chunk.Content) > 0 {
			if !content {
				fmt.Printf("\n\n--- Content ---\n")
				content = true
				reasoning = false
				toolCalls = false
			}
			fmt.Print(chunk.Content)
		}
		if len(chunk.ToolCalls) > 0 {
			if !toolCalls {
				fmt.Printf("\n\n--- Tool Calls ---\n")
				toolCalls = true
				reasoning = false
				content = false
			}
			fmt.Println(chunk.ToolCalls)
		}
	}
	fmt.Print("\n\n")

	fmt.Println(runResponse.Content())

	defer logger.SaveToFile("!!_requests.json")
}
