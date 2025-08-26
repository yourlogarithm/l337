package main

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/provider/ollama"
	"github.com/yourlogarithm/l337/run"
)

func main() {
	client, logger := newLoggingHTTPClient()

	// model := openai.NewModel(
	// 	"zai-org/GLM-4.5",
	// 	option.WithBaseURL(os.Getenv("BASE_URL")),
	// 	option.WithAPIKey(os.Getenv("API_KEY")),
	// 	option.WithHTTPClient(client),
	// )

	model, _ := ollama.NewModel(
		"qwen3:8b",
		"http://localhost:11434",
		// os.Getenv("OLLAMA_BASE_URL"),
		client,
	)

	chatOptions := provider.ChatOptions{
		ReasoningEffort: provider.NewReasoningEffortBool(true),
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

	response, err := team.RunWithParams(
		context.Background(),
		run.WithMessage(chat.RoleUser, "Discuss the pros and cons of using AI in education."),
	)
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}

	// fmt.Println("Team Response:", response.Content())
	fmt.Printf("%v\n", response.Metrics)

	defer logger.SaveToFile("requests.json")
}
