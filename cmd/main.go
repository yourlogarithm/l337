package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
)

func main() {
	client, logger := newLoggingHTTPClient()

	model := openai.NewModel(
		"Qwen/Qwen3-0.6B",
		option.WithBaseURL(os.Getenv("BASE_URL")),
		option.WithAPIKey(os.Getenv("API_KEY")),
		option.WithHTTPClient(client),
	)

	// model, _ := ollama.NewModel(
	// 	"qwen3:8b",
	// 	"http://localhost:11434",
	// 	// os.Getenv("OLLAMA_BASE_URL"),
	// 	client,
	// )

	chatOptions := chat.Options{
		// ReasoningEffort: chat.NewReasoningEffortBool(true),
		IncludeStreamMetrics: true,
	}

	inFavorAgent, err := agent.New(model, agent.WithName("favor_agent"), agent.WithDescription("Agent that provides a detailed analysis in favor of the discussed topic."), agent.WithInstructions("Provide strong arguments and detailed analysis. Use point-by-point structure. Respond with a single side of the argument in markdown format."), agent.WithChatOptions(chatOptions))
	if err != nil {
		panic(err)
	}

	againstAgent, err := agent.New(model, agent.WithName("against_agent"), agent.WithDescription("Agent that provides a detailed analysis against the discussed topic."), agent.WithInstructions("Provide strong arguments and detailed analysis. Use point-by-point structure. Respond with a single side of the argument in markdown format."), agent.WithChatOptions(chatOptions))
	if err != nil {
		panic(err)
	}

	team, err := agent.New(model, agent.WithName("debate_team"), agent.WithDescription("A team of agents debating a topic."), agent.WithInstructions("Use your team members to collaboratively analyze the topic and provide a comprehensive response.\nUse InFavorAgent to assign him the task of providing a supportive perspective on the topic.\nUse AgainstAgent to assign him the task of providing an opposing perspective on the topic.\nAfter analyzing both responses, you must come to a conclusion choosing a single side which had better arguments."), agent.WithSubordinate(inFavorAgent), agent.WithSubordinate(againstAgent), agent.WithChatOptions(chatOptions))
	if err != nil {
		panic(err)
	}

	runResponse := chat.RunResponse{
		Messages: []chat.Message{
			{
				Role:    chat.RoleUser,
				Content: "Discuss the pros and cons of using AI in education.",
			},
		},
	}

	stream, err := team.RunStreaming(
		context.Background(),
		&runResponse,
		128,
	)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}

	fmt.Println("Consuming stream...")
	reasoning := false
	for {
		chunk, err := stream.Next()
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

	metrics, _ := json.MarshalIndent(runResponse.Metrics, "", "  ")
	fmt.Printf("\n%s\n", metrics)

	defer logger.SaveToFile("requests.json")
}
