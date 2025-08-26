package main

import (
	"net/http"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/agent/remote"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/provider/ollama"
)

func main() {
	model, _ := ollama.NewModel(
		"qwen3:8b",
		"http://localhost:11434",
		http.DefaultClient,
	)
	chatOptions := provider.ChatOptions{
		ReasoningEffort: provider.NewReasoningEffortBool(true),
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
	server.Serve(":8080", nil)
}
