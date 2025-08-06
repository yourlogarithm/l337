package main

import (
	"context"
	"fmt"
	"os"

	"github.com/yourlogarithm/golagno/agent"
	"github.com/yourlogarithm/golagno/agentic"
	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/provider"
	"github.com/yourlogarithm/golagno/provider/ollama"
	"github.com/yourlogarithm/golagno/team"
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

// Returns the difference of a and b
func subtract(ctx context.Context, params tools.Params) (string, error) {
	a, err := tools.GetParameter[int](params, "a")
	if err != nil {
		return "", err
	}
	b, err := tools.GetParameter[int](params, "b")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", a-b), nil
}

func math(model *provider.Model) {
	toolkit := tools.Toolkit{}

	addTool := tools.NewTool("add", "Adds two numbers", add)
	tools.AddParameterFromType[int](&addTool, "a", "The first number to add", true)
	tools.AddParameterFromType[int](&addTool, "b", "The second number to add", true)

	subtractTool := tools.NewTool("subtract", "Subtracts two numbers", subtract)
	tools.AddParameterFromType[int](&subtractTool, "a", "The first number to subtract", true)
	tools.AddParameterFromType[int](&subtractTool, "b", "The second number to subtract", true)

	toolkit.AddTool(addTool)
	toolkit.AddTool(subtractTool)

	options := agentic.Options{
		Name:         "ExampleAgent",
		Description:  "An example agent for demonstration purposes.",
		Instructions: "None at all",
		Model:        model,
		Tools:        toolkit,
	}
	agent := agent.NewFromOptions(options)
	response, err := agent.Run(context.Background(), []chat.Message{
		{
			Role:    "user",
			Content: "Perform the following operations: 5 + 3, 23 + 42, 66 - 33",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Agent Response:", response.Content())
}

func collaborate(model *provider.Model) {
	aOptions := agentic.Options{
		Name:         "InFavorAgent",
		Description:  "Agent that provides a detailed analysis in favor of the discussed topic.",
		Instructions: "Provide strong arguments and detailed analysis. Use point-by-point structure. Respond with a single side of the argument in markdown format.",
		Model:        model,
	}
	inFavorAgent := agent.NewFromOptions(aOptions)

	bOptions := agentic.Options{
		Name:         "AgainstAgent",
		Description:  "Agent that provides a detailed analysis against the discussed topic.",
		Instructions: "Provide strong arguments and detailed analysis. Use point-by-point structure. Respond with a single side of the argument in markdown format.",
		Model:        model,
	}
	againstAgent := agent.NewFromOptions(bOptions)

	teamOptions := agentic.Options{
		Name:         "DebateTeam",
		Description:  "A team of agents debating a topic.",
		Instructions: "Use your team members to collaboratively analyze the topic and provide a comprehensive response.\nUse InFavorAgent to assign him the task of providing a supportive perspective on the topic.\nUse AgainstAgent to assign him the task of providing an opposing perspective on the topic.\nAfter analyzing both responses, you must come to a conclusion choosing a single side which had better arguments.",
		Model:        model,
	}

	t := team.Team{
		Options: teamOptions,
		Members: []agentic.Member{
			inFavorAgent,
			againstAgent,
		},
		Mode: team.ModeCollaborate,
	}

	response, err := t.Run(context.Background(), []chat.Message{
		{
			Role:    chat.RoleUser.String(),
			Content: "Discuss the pros and cons of using AI in education.",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Team Response:", response.Content())
}

func main() {

	client, logger := newLoggingHTTPClient()

	// model := openai.NewOpenAI(
	// 	"Qwen/QwQ-32B",
	// 	option.WithBaseURL(os.Getenv("BASE_URL")),
	// 	option.WithAPIKey(os.Getenv("API_KEY")),
	// 	option.WithHTTPClient(client),
	// )

	model, _ := ollama.NewOllama(
		"gpt-oss:20b",
		os.Getenv("OLLAMA_BASE_URL"),
		client,
	)

	collaborate(model)

	defer logger.SaveToFile("requests.json")
}
