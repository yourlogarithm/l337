package main

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/provider/ollama"
)

// Returns the sum of a and b
// func add(ctx context.Context, params tools.Params) (string, error) {
// 	a, err := tools.GetParameter[int](&params, "a")
// 	if err != nil {
// 		return "", err
// 	}
// 	b, err := tools.GetParameter[int](&params, "b")
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprintf("%d", a+b), nil
// }

// Returns the difference of a and b
// func subtract(ctx context.Context, params tools.Params) (string, error) {
// 	a, err := tools.GetParameter[int](&params, "a")
// 	if err != nil {
// 		return "", err
// 	}
// 	b, err := tools.GetParameter[int](&params, "b")
// 	if err != nil {
// 		return "", err
// 	}
// 	return fmt.Sprintf("%d", a-b), nil
// }

// func math(model *provider.Model) {
// 	addTool := tools.NewTool("add", "Adds two numbers", add)
// 	tools.AddParameterFromType[int](&addTool, "a", "The first number to add", true)
// 	tools.AddParameterFromType[int](&addTool, "b", "The second number to add", true)

// 	subtractTool := tools.NewTool("subtract", "Subtracts two numbers", subtract)
// 	tools.AddParameterFromType[int](&subtractTool, "a", "The first number to subtract", true)
// 	tools.AddParameterFromType[int](&subtractTool, "b", "The second number to subtract", true)

// 	agent, err := agent.New(model, agent.WithName("example_agent"), agent.WithDescription("An example agent for demonstration purposes."), agent.WithInstructions("None at all"), agent.WithTool(addTool), agent.WithTool(subtractTool))
// 	if err != nil {
// 		panic(err)
// 	}
// 	response, err := agent.Run(context.Background(), []chat.Message{
// 		{
// 			Role:    "user",
// 			Content: "Perform the following operations: 5 + 3, 23 + 42, 66 - 33",
// 		},
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Agent Response:", response.Content())
// }

func collaborate(model *provider.Model) {
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

	response, err := team.Run(context.Background(), []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Discuss the pros and cons of using AI in education.",
		},
	})
	if err != nil {
		fmt.Printf("Error occurred: %v\n", err)
	}
	fmt.Println("Team Response:", response.Content())
}

func main() {

	// type a struct {
	// 	X int                 `json:"x"`
	// 	Y []float64           `json:"y"`
	// 	Z map[string][]string `json:"z"`
	// }

	// type b struct {
	// 	A a `json:"a"`
	// }

	// schema := jsonschema.Reflect(b{})

	// targetRef := strings.TrimPrefix(schema.Ref, "#/$defs/")
	// v, ok := schema.Definitions[targetRef]
	// if !ok {
	// 	panic(fmt.Sprintf("definition %s not found", targetRef))
	// }
	// schema.Properties = v.Properties
	// delete(schema.Definitions, targetRef)

	// marshaled, _ := json.MarshalIndent(schema, "", "  ")
	// fmt.Println(string(marshaled))

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

	collaborate(model)

	defer logger.SaveToFile("requests.json")
}
