package examples

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

type AddParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func add(ctx context.Context, response *run.Response, addParams AddParams) (string, error) {
	return fmt.Sprintf("%f", addParams.A+addParams.B), nil
}

type SubtractParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func subtract(ctx context.Context, response *run.Response, subtractParams SubtractParams) (string, error) {
	return fmt.Sprintf("%f", subtractParams.A-subtractParams.B), nil
}

type MultiplyParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func multiply(ctx context.Context, response *run.Response, multiplyParams MultiplyParams) (string, error) {
	return fmt.Sprintf("%f", multiplyParams.A*multiplyParams.B), nil
}

type DivideParams struct {
	A float32 `json:"a" jsonschema:"required"`
	B float32 `json:"b" jsonschema:"required"`
}

func divide(ctx context.Context, response *run.Response, divideParams DivideParams) (string, error) {
	if divideParams.B != 0 {
		return fmt.Sprintf("%f", divideParams.A/divideParams.B), nil
	}
	return "division by zero error", nil
}

func ToolsExample() {
	model := openai.NewModel("gpt-4o")

	addTool, _ := tools.NewToolWithArgs("add", "Adds two numbers", add)
	subtractTool, _ := tools.NewToolWithArgs("subtract", "Subtracts two numbers", subtract)
	multiplyTool, _ := tools.NewToolWithArgs("multiply", "Multiplies two numbers", multiply)
	divideTool, _ := tools.NewToolWithArgs("divide", "Divides two numbers", divide)

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
		run.WithMessage(chat.RoleUser, "What is 5 + 3?"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())
}
