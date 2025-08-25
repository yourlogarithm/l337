package examples

import (
	"context"
	"fmt"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
	"github.com/yourlogarithm/l337/run"
)

func TeamExample() {
	model := openai.NewModel("gpt-4o")

	jediAgent, _ := agent.New(
		model,
		agent.WithName("jedi_master"),
		agent.WithDescription("A wise and powerful Jedi Master."),
		agent.WithInstructions(
			"You are a wise and powerful Jedi Master. "+
				"Use your knowledge of the Force to assist the user. "+
				"Be calm, patient, and wise in your responses. "+
				"Respond to the user's queries like a Jedi Master would.",
		),
	)

	sithAgent, _ := agent.New(
		model,
		agent.WithName("sith_lord"),
		agent.WithDescription("A cunning and powerful Sith Lord."),
		agent.WithInstructions(
			"You are a cunning and powerful Sith Lord. "+
				"Use your knowledge of the Dark Side to assist the user. "+
				"Be aggressive, cunning, and ruthless in your responses. "+
				"Respond to the user's queries like a Sith Lord would.",
		),
	)

	starWarsTeam, _ := agent.New(
		model,
		agent.WithName("star_wars_team"),
		agent.WithDescription("A team of agents representing the Jedi and Sith."),
		agent.WithInstructions(
			"You are a team of agents representing the Jedi and Sith. "+
				"Choose the best agent to respond to the user's queries based on their nature. "+
				"The Jedi will provide answers about their ways, while the Sith will provide answers about the Dark Side. "+
				"If the user asks a question that is unrelated to either, respond directly that you cannot answer that question. "+
				"If the question is related to Star Wars lore, but is neutral, respond directly without assigning it to either agent.",
		),
		agent.WithSubordinate(sithAgent),
		agent.WithSubordinate(jediAgent),
	)

	response, err := starWarsTeam.RunWithParams(
		context.Background(),
		run.WithMessage(chat.RoleUser, "Who was the most powerful Jedi?"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())

	response, err = starWarsTeam.RunWithParams(
		context.Background(),
		run.WithMessage(chat.RoleUser, "Who was the most powerful Sith?"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())
}
