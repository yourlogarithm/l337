package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/yourlogarithm/golagno/agent"
	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/provider"
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

	slog.SetLogLoggerLevel(slog.LevelDebug)

	base, _ := url.Parse("http://localhost:11434")
	model := provider.NewOllama(
		"llama3.2:1b",
		base,
		http.DefaultClient,
	)
	toolkit := tools.Toolkit{}

	tool := tools.NewTool("add", "Adds two numbers", add)
	tools.AddParameterFromType[int](&tool, "a")
	tools.AddParameterFromType[int](&tool, "b")

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
	fmt.Println("Agent Response:", response)
}
