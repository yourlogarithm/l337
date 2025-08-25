package agent

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/retry"
	"github.com/yourlogarithm/l337/tools"
)

var logger = logging.SetupLogger("agent")

type Agent struct {
	// Unique identifier for the member
	// If not set, a random UUID will be generated
	id uuid.UUID
	// Appended to the system message
	// If not set, an error will be returned as soon as `Initialize` is called
	name string
	// Provided to the parent member (if any) to help understand the member's purpose
	role string
	// Description of the member
	// Appends to the system message
	description string
	// Appended to the system message
	// Wrapped in <instructions> tags
	instructions string
	// Appended to the system message
	// Wrapped in <goal> tags
	goal string
	// Appended to the system message
	// Wrapped in <expected_output> tags
	expectedOutput string
	// Model used to send LLM requests
	model *provider.Model
	// Tools for the LLM to use
	tools tools.Toolkit
	// Retry options for the LLM requests
	// If not set, defaults to `retry.Default()`
	retry *retry.Options
	// List of subordinate agents that this agent can delegate tasks to
	subordinates []AgentImpl
	chatOptions  provider.ChatOptions
}

func (a *Agent) Name() string {
	return a.name
}

func (a *Agent) Description() string {
	return a.description
}

func (a *Agent) Skills() []tools.SkillCard {
	var skills = make([]tools.SkillCard, 0, len(a.tools))
	for _, tool := range a.tools {
		skills = append(skills, tool.SkillCard)
	}
	return skills
}
