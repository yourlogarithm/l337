package main

import (
	"fmt"

	"github.com/yourlogarithm/golagno/agent"
	"github.com/yourlogarithm/golagno/provider"
)

func main() {
	model := provider.NewOpenAI("gpt-3.5-turbo")
	agent := agent.Agent{
		Name:         "ExampleAgent",
		Description:  "An example agent for demonstration purposes.",
		Instructions: "None at all",
		Model:        model,
	}
	response, err := agent.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println("Agent Response:", response)
}
