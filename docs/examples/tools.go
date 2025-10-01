package examples

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
	"github.com/yourlogarithm/l337/tools"
)

type AddParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func add(ctx context.Context, response *chat.RunResponse, addParams AddParams) ([]chat.Content, error) {
	return chat.NewTextContent(fmt.Sprintf("%f", addParams.A+addParams.B)).AsSlice(), nil
}

type SubtractParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func subtract(ctx context.Context, response *chat.RunResponse, subtractParams SubtractParams) ([]chat.Content, error) {
	return chat.NewTextContent(fmt.Sprintf("%f", subtractParams.A-subtractParams.B)).AsSlice(), nil
}

type MultiplyParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func multiply(ctx context.Context, response *chat.RunResponse, multiplyParams MultiplyParams) ([]chat.Content, error) {
	return chat.NewTextContent(fmt.Sprintf("%f", multiplyParams.A*multiplyParams.B)).AsSlice(), nil
}

type DivideParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func divide(ctx context.Context, response *chat.RunResponse, divideParams DivideParams) ([]chat.Content, error) {
	if divideParams.B != 0 {
		return chat.NewTextContent(fmt.Sprintf("%f", divideParams.A/divideParams.B)).AsSlice(), nil
	}
	return chat.NewTextContent("division by zero error").AsSlice(), nil
}

func ToolsExample() {
	model := openai.NewModel("gpt-4o")

	addTool, _ := tools.NewWithArgs("add", "Adds two numbers", add)
	subtractTool, _ := tools.NewWithArgs("subtract", "Subtracts two numbers", subtract)
	multiplyTool, _ := tools.NewWithArgs("multiply", "Multiplies two numbers", multiply)
	divideTool, _ := tools.NewWithArgs("divide", "Divides two numbers", divide)

	mathAgent, err := agent.New(
		model,
		agent.WithName("math_agent"),
		agent.WithDescription("An agent that can perform basic math operations."),
		agent.WithInstructions("Perform addition, subtraction, multiplication, and division based on the user's request. If the operation is not supported, respond that you cannot perform that operation."),
		agent.WithTool(addTool),
		agent.WithTool(subtractTool),
		agent.WithTool(multiplyTool),
		agent.WithTool(divideTool),
	)
	if err != nil {
		panic(err)
	}

	response, err := mathAgent.RunWithParams(
		context.Background(),
		chat.WithTextMessage(chat.RoleUser, "What is 5 + 3?"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())
}
