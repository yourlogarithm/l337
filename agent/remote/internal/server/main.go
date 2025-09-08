package main

import (
	"fmt"
	"os"

	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/agent/remote"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
)

func main() {
	// model, _ := ollama.NewModel(
	// 	"qwen3:8b",
	// 	"http://localhost:11434",
	// 	http.DefaultClient,
	// )
	model := openai.NewModel(
		"Qwen/Qwen3-0.6B",
		option.WithBaseURL(os.Getenv("BASE_URL")),
		option.WithAPIKey(os.Getenv("API_KEY")),
		// option.WithHTTPClient(client),
	)
	chatOptions := chat.Options{
		// 	ReasoningEffort: chat.NewReasoningEffortBool(true),
		IncludeStreamMetrics: true,
	}
	agent, _ := agent.New(
		model,
		agent.WithName("remote_agent"),
		agent.WithDescription("An agent accessible via remote calls."),
		agent.WithInstructions("You are a helpful assistant."),
		agent.WithChatOptions(chatOptions),
	)
	server := remote.AgentServer{
		Agent: agent,
	}

	fmt.Println("Starting remote agent server on :8080")
	server.Serve(":8080", nil)
}
