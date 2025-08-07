<table>
    <tr>
        <td>
            <img src="docs/logo.svg" alt="Logo" width="150px">
        </td>
        <td>
            <h1>l337</h1>
            <p>l337 (pronounced "leet") is a framework inspired by <a href="https://github.com/agno-agi/agno">Agno</a> for building multi-agent systems.</p>
        </td>
    </tr>
</table>

# Installation
```bash
go get github.com/yourlogarithm/l337
```

# Usage

## Agent

```go
import (
    "fmt"
    "context"
    "github.com/yourlogarithm/l337/agent"
    "github.com/yourlogarithm/l337/agentic"
    "github.com/yourlogarithm/l337/chat"
    "github.com/yourlogarithm/l337/provider/openai"
)

func main() {
    model := openai.NewModel("gpt-4o")
    opts := agentic.Options{
        Name: "Obi-Wan Kenobi",
        Description: "A wise and powerful Jedi Master.",
        Instructions: "You are Obi-Wan Kenobi, a wise and powerful Jedi Master. Use your knowledge of the Force to assist the user. Be calm, patient, and wise in your responses. Respond to the user's queries like a Obi-Wan Kenobi would.",
        Model: model,
    }
    agent, err := agent.NewFromOptions(opts)
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
```

## Team
```go 
import (
	"context"
	"fmt"
	"os"

	"github.com/openai/openai-go/option"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider/openai"
	"github.com/yourlogarithm/l337/team"
)

func main() {
	model := openai.NewModel(
		"Qwen/QwQ-32B",
		option.WithBaseURL(os.Getenv("BASE_URL")),
		option.WithAPIKey(os.Getenv("API_KEY")),
	)

	jediOptions := agentic.Options{
		Name:         "Jedi Master",
		Role:         "Help the user with questions regarding the Jedi ways.",
		Description:  "A wise and powerful Jedi Master.",
		Instructions: "You are a wise and powerful Jedi Master. Use your knowledge of the Force to assist the user. Be calm, patient, and wise in your responses. Respond to the user's queries like a Jedi Master would.",
		Model:        model,
	}
	jediAgent, err := agent.NewFromOptions(jediOptions)
	if err != nil {
		panic(err)
	}

	sithOptions := agentic.Options{
		Name:         "Sith Lord",
		Role:         "Help the user with questions regarding the Dark Side of the Force.",
		Description:  "A cunning and powerful Sith Lord.",
		Instructions: "You are a cunning and powerful Sith Lord. Use your knowledge of the Dark Side to assist the user. Be aggressive, cunning, and ruthless in your responses. Respond to the user's queries like a Sith Lord would.",
		Model:        model,
	}
	sithAgent, err := agent.NewFromOptions(sithOptions)
	if err != nil {
		panic(err)
	}

	teamOptions := agentic.Options{
		Name:         "Star Wars Team",
		Description:  "A team of agents representing the Jedi and Sith.",
		Instructions: "You are a team of agents representing the Jedi and Sith. Choose the best agent to respond to the user's queries based on their nature. The Jedi will provide answers about their ways, while the Sith will provide answers about the Dark Side. If the user asks a question that is unrelated to either, respond directly that you cannot answer that question. If the question is related to Star Wars lore, but is neutral, respond directly without assigning it to either agent.",
		Model:        model,
	}
	team := team.Team{
		Options: teamOptions,
		Members: []agentic.Member{
			jediAgent,
			sithAgent,
		},
		Mode: team.ModeRoute,
	}

	messages := []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Who was the most powerful Jedi?",
		},
	}
	response, err := team.Run(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())

	messages = []chat.Message{
		{
			Role:    chat.RoleUser,
			Content: "Who was the most powerful Sith?",
		},
	}
	response, err = team.Run(context.Background(), messages)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Content())
}
```