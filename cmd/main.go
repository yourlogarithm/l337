package main

import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/provider/openai"
	"github.com/yourlogarithm/l337/team"
	"github.com/yourlogarithm/l337/tools"
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
	agent, err := agent.NewFromOptions(options)
	if err != nil {
		panic(err)
	}
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
	inFavorAgent, err := agent.NewFromOptions(aOptions)
	if err != nil {
		panic(err)
	}

	bOptions := agentic.Options{
		Name:         "AgainstAgent",
		Description:  "Agent that provides a detailed analysis against the discussed topic.",
		Instructions: "Provide strong arguments and detailed analysis. Use point-by-point structure. Respond with a single side of the argument in markdown format.",
		Model:        model,
	}
	againstAgent, err := agent.NewFromOptions(bOptions)
	if err != nil {
		panic(err)
	}

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
			Role:    chat.RoleUser,
			Content: "Discuss the pros and cons of using AI in education.",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Team Response:", response.Content())
}

func route(model *provider.Model) {
	englishAgent, err := agent.NewFromOptions(agentic.Options{
		Name:         "English Agent",
		Description:  "An agent that specializes in English language tasks.",
		Instructions: "You must only respond in English",
		Model:        model,
	})
	if err != nil {
		panic(err)
	}

	chineseAgent, err := agent.NewFromOptions(agentic.Options{
		Name:         "Chinese Agent",
		Description:  "An agent that specializes in Chinese language tasks.",
		Instructions: "You must only respond in Chinese",
		Model:        model,
	})
	if err != nil {
		panic(err)
	}

	frenchAgent, err := agent.NewFromOptions(agentic.Options{
		Name:         "French Agent",
		Description:  "An agent that specializes in French language tasks.",
		Instructions: "You must only respond in French",
		Model:        model,
	})
	if err != nil {
		panic(err)
	}

	multiLanguageTeam := team.Team{
		Options: agentic.Options{
			Name:         "Multi-Language Team",
			Description:  "A team of agents that can handle tasks in multiple languages.",
			Instructions: "You are a language router that directs questions to the appropriate language agent.\nIf the user asks in a language whose agent is not a team member, respond in English with:\n'I can only answer in the following languages: English, Spanish, Japanese, French and German. Please ask your question in one of these languages.'\nAlways check the language of the user's input before routing to an agent.\nFor unsupported languages like Italian, respond in English with the above message.",
			Model:        model,
		},
		Members: []agentic.Member{
			englishAgent,
			chineseAgent,
			frenchAgent,
		},
		Mode: team.ModeRoute,
	}

	messages := []string{
		"Hello, how are you?",
		"你好吗？",
		"Bonjour, comment ça va ?",
		"¿Cómo estás?",
		"こんにちは、お元気ですか？",
	}
	ctx := context.Background()
	for _, msg := range messages {
		response, err := multiLanguageTeam.Run(ctx, []chat.Message{
			{
				Role:    chat.RoleUser,
				Content: msg,
			},
		})
		if err != nil {
			fmt.Printf("Error processing message '%s': %v\n", msg, err)
			continue
		}
		fmt.Printf("Response to '%s': %s\n", msg, response.Content())
	}
}

func main() {

	client, logger := newLoggingHTTPClient()

	model := openai.NewOpenAI(
		"Qwen/QwQ-32B",
		option.WithBaseURL(os.Getenv("BASE_URL")),
		option.WithAPIKey(os.Getenv("API_KEY")),
		option.WithHTTPClient(client),
	)

	// model, _ := ollama.NewOllama(
	// 	"gpt-oss:20b",
	// 	os.Getenv("OLLAMA_BASE_URL"),
	// 	client,
	// )

	route(model)

	defer logger.SaveToFile("requests.json")
}
