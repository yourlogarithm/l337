package examples

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
	"github.com/yourlogarithm/l337/tools"
)

func add(ctx context.Context, params tools.Params) (string, error) {
	a, err := tools.GetParameter[float32](params, "a")
	if err != nil {
		return "", err
	}
	b, err := tools.GetParameter[float32](params, "b")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%f", a+b), nil
}

func subtract(ctx context.Context, params tools.Params) (string, error) {
	a, err := tools.GetParameter[float32](params, "a")
	if err != nil {
		return "", err
	}
	b, err := tools.GetParameter[float32](params, "b")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%f", a-b), nil
}

func multiply(ctx context.Context, params tools.Params) (string, error) {
	a, err := tools.GetParameter[float32](params, "a")
	if err != nil {
		return "", err
	}
	b, err := tools.GetParameter[float32](params, "b")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%f", a*b), nil
}

func divide(ctx context.Context, params tools.Params) (string, error) {
	a, err := tools.GetParameter[float32](params, "a")
	if err != nil {
		return "", err
	}
	b, err := tools.GetParameter[float32](params, "b")
	if err != nil {
		return "", err
	}
	if b != 0 {
		return fmt.Sprintf("%f", a/b), nil
	}
	return "division by zero error", nil
}

func ToolsExample() {
	model := openai.NewModel("gpt-4o")

	mathAgentOptions := agentic.Configuration{
		Name:         "Math Agent",
		Description:  "An agent that can perform basic math operations.",
		Instructions: "Perform addition, subtraction, multiplication, and division based on the user's request. If the operation is not supported, respond that you cannot perform that operation.",
		Model:        model,
	}

	addTool := tools.NewTool("add", "Adds two numbers", add)
	tools.AddParameterFromType[float32](&addTool, "a", "The first number to add", true)
	tools.AddParameterFromType[float32](&addTool, "b", "The second number to add", true)
	mathAgentOptions.Tools.AddTool(addTool)

	subtractTool := tools.NewTool("subtract", "Subtracts two numbers", subtract)
	tools.AddParameterFromType[float32](&subtractTool, "a", "The first number to subtract", true)
	tools.AddParameterFromType[float32](&subtractTool, "b", "The second number to subtract", true)
	mathAgentOptions.Tools.AddTool(subtractTool)

	multiplyTool := tools.NewTool("multiply", "Multiplies two numbers", multiply)
	tools.AddParameterFromType[float32](&multiplyTool, "a", "The first number to multiply", true)
	tools.AddParameterFromType[float32](&multiplyTool, "b", "The second number to multiply", true)
	mathAgentOptions.Tools.AddTool(multiplyTool)

	divideTool := tools.NewTool("divide", "Divides two numbers", divide)
	tools.AddParameterFromType[float32](&divideTool, "a", "The first number to divide", true)
	tools.AddParameterFromType[float32](&divideTool, "b", "The second number to divide", true)
	mathAgentOptions.Tools.AddTool(divideTool)

	mathAgent, err := agent.NewFromOptions(mathAgentOptions)
	if err != nil {
		panic(err)
	}

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "What is 5 + 3?",
		},
	}
	response, err := mathAgent.Run(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())
}
