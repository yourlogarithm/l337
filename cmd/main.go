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
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	base, _ := url.Parse("http://localhost:11434")
	model := provider.NewOllama(
		"gemma3:4b",
		base,
		http.DefaultClient,
	)
	agent := agent.Agent{
		Name:         "ExampleAgent",
		Description:  "An example agent for demonstration purposes.",
		Instructions: "None at all",
		Model:        model,
	}
	response, err := agent.Run(context.Background(), []chat.Message{
		{
			Role:    "user",
			Content: "Hello, how are you?",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Agent Response:", response)
}
