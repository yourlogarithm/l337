package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/agent/remote"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/provider/ollama"
	"github.com/yourlogarithm/l337/run"
)

func main() {
	remoteAgent := remote.DefaultClient("http://localhost:8080")

	model, _ := ollama.NewModel(
		"qwen3:8b",
		"http://localhost:11434",
		http.DefaultClient,
	)
	chatOptions := provider.ChatOptions{
		ReasoningEffort: provider.NewReasoningEffortBool(true),
	}
	agentWrapper, _ := agent.New(
		model,
		agent.WithName("wrapper_agent"),
		agent.WithInstructions("Route all messages to the subordinate agent"),
		agent.WithDescription("A wrapper agent that delegates tasks to a subordinate agent"),
		agent.WithChatOptions(chatOptions),
		agent.WithSubordinate(remoteAgent),
	)

	response, err := agentWrapper.RunWithParams(
		context.Background(),
		run.WithMessage(chat.RoleUser, "Hello, ask your subordinate what time it is."),
		run.WithSessionID(uuid.New()),
	)
	if err != nil {
		panic(err)
	}

	marshaled, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(marshaled))
}
