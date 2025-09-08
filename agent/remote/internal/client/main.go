package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/yourlogarithm/l337/agent/remote/client"
	"github.com/yourlogarithm/l337/chat"
)

func main() {
	remoteAgent := client.Default("http://localhost:8080")

	// model := openai.NewModel(
	// 	"Qwen/Qwen3-0.6B",
	// 	option.WithBaseURL(os.Getenv("BASE_URL")),
	// 	option.WithAPIKey(os.Getenv("API_KEY")),
	// )
	// agentWrapper, _ := agent.New(
	// 	model,
	// 	agent.WithName("wrapper_agent"),
	// 	agent.WithInstructions("Route all messages to the subordinate agent"),
	// 	agent.WithDescription("A wrapper agent that delegates tasks to a subordinate agent"),
	// 	agent.WithSubordinate(remoteAgent),
	// )

	// response, err := agentWrapper.RunWithParams(
	// 	context.Background(),
	// 	chat.WithMessage(chat.RoleUser, "Hello, ask your subordinate what time it is."),
	// 	chat.WithSessionID(uuid.New()),
	// )
	// if err != nil {
	// 	panic(err)
	// }

	// marshaled, _ := json.MarshalIndent(response, "", "  ")
	// os.WriteFile("conversation.json", marshaled, 0644)

	runResponse, ch, err := remoteAgent.RunStreamingWithParams(
		context.Background(),
		128,
		chat.WithMessage(chat.RoleUser, "Write a funny joke."),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consuming stream...")
	reasoning := false
	for {
		chunk, err := ch.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("Error occurred: %v\n", err)
			break
		}
		if chunk.Content != "" {
			if reasoning {
				reasoning = false
				fmt.Printf("\n\n--- Content ---\n")
			}
			fmt.Printf("%s", chunk.Content)
		} else if chunk.Reasoning != "" {
			if !reasoning {
				fmt.Printf("\n\n--- Reasoning ---\n")
				reasoning = true
			}
			fmt.Printf("%s", chunk.Reasoning)
		}
	}

	fmt.Println()

	marshaled, _ := json.MarshalIndent(runResponse, "", "  ")
	fmt.Println(string(marshaled))
}
