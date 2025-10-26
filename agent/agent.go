package agent

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/providers"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

var logger = logging.SetupLogger("agent")

type Agent struct {
	// Unique identifier for the member
	// If not set, a random UUID will be generated
	id uuid.UUID
	// Appended to the system message
	// If not set, an error will be returned as soon as `Initialize` is called
	name string
	// Description of the member
	// Appends to the system message as is, with a newline right after the name
	description string
	// Appended to the system message
	// Wrapped in <instructions> tags
	instructions string
	// Appended to the system message
	// Wrapped in <expected_output> tags
	expectedOutput string
	// Model used to send LLM requests
	model *providers.Model
	// Tools for the LLM to use
	tools tools.Toolkit
	// Retry options for the LLM requests
	// If not set, defaults to `retry.Default()`
	// retry *retry.Options
	// List of subordinate agents that this agent can delegate tasks to
	subordinates []AgentImpl
	chatOptions  types.Options
}

func (a *Agent) Name() (string, error) {
	return a.name, nil
}

func (a *Agent) Description() (string, error) {
	return a.description, nil
}

func (a *Agent) Tools() ([]tools.Tool, error) {
	var tools = make([]tools.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, tool)
	}
	return tools, nil
}
