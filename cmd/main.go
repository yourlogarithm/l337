package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/golagno/agent"
	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/provider/openai"
	"github.com/yourlogarithm/golagno/tools"
)

// Returns the sum of a and b
func add(ctx context.Context, params tools.Params) (string, error) {
	a, err := tools.GetParameter[int](params, "a")
	if err != nil {
		return "", err
	}
	b, err := tools.GetParameter[int](params, "b")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", a+b), nil
}

func main() {
	model := openai.NewOpenAI(
		"Qwen/QwQ-32B",
		option.WithBaseURL(os.Getenv("BASE_URL")),
		option.WithAPIKey(os.Getenv("API_KEY")),
	)
	toolkit := tools.Toolkit{}

	tool := tools.NewTool("add", "Adds two numbers", add)
	tools.AddParameterFromType[int](&tool, "a", "The first number to add")
	tools.AddParameterFromType[int](&tool, "b", "The second number to add")

	toolkit.AddTool(tool)
	agent := agent.Agent{
		Name:         "ExampleAgent",
		Description:  "An example agent for demonstration purposes.",
		Instructions: "None at all",
		Model:        model,
		Tools:        toolkit,
	}
	response, err := agent.Run(context.Background(), []chat.Message{
		{
			Role:    "user",
			Content: "Perform the following operations: 5 + 3, 23 + 42, 66 + 33",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Agent Response:", response.Content())
}
