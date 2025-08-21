package examples

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
)

func AgentExample() {
	model := openai.NewModel("gpt-4o")
	agent, err := agent.New(
		model,
		agent.WithName("obi_wan_kenobi"),
		agent.WithDescription("A wise and powerful Jedi Master."),
		agent.WithInstructions(
			"You are Obi-Wan Kenobi, a wise and powerful Jedi Master. "+
				"Use your knowledge of the Force to assist the user. "+
				"Be calm, patient, and wise in your responses. "+
				"Respond to the user's queries like a Obi-Wan Kenobi would.",
		),
	)
	if err != nil {
		panic(err)
	}

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "What was your greatest failure, and what did you learn from it?",
		},
	}
	response, err := agent.Run(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())
}
