package main

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/ollama"
	"github.com/yourlogarithm/l337/tools"
)

type FooParams struct {
	I int `json:"superMegaUltra"`
}

func foo(ctx context.Context, runResponse *chat.RunResponse, params FooParams) (string, error) {
	return strconv.Itoa(params.I * 3), nil
}

func main() {
	client, logger := newLoggingHTTPClient()

	// model := openai.NewModel(
	// 	"gpt-5-nano",
	// 	// option.WithBaseURL(os.Getenv("BASE_URL")),
	// 	option.WithAPIKey(os.Getenv("API_KEY")),
	// 	option.WithHTTPClient(client),
	// )

	model, _ := ollama.NewModel(
		"qwen3:8b",
		"http://localhost:11434",
		client,
	)

	tool, err := tools.NewWithArgs("foo", "Foo tool", foo)
	if err != nil {
		panic(err)
	}

	options := chat.Options{
		// ReasoningEffort: chat.NewReasoningEffortBool(true),
		ReasoningEffort: chat.NewReasoningEffortBool(true),
	}

	a, err := agent.New(model, agent.WithTool(tool), agent.WithChatOptions(options))
	if err != nil {
		panic(err)
	}

	runResponse, stream, err := a.RunStreamingWithParams(
		context.Background(),
		0,
		chat.WithMessage(chat.RoleUser, "Concurrently call `foo` tool 3 times with a different value."),
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
