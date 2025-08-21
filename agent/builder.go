package agent

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/retry"
	"github.com/yourlogarithm/l337/tools"
)

type AgentOption interface {
	Apply(*Agent) error
}

type AgentOptionFunc func(*Agent) error

func (s AgentOptionFunc) Apply(r *Agent) error { return s(r) }

func New(model *provider.Model, options ...AgentOption) (*Agent, error) {
	if model == nil {
		return nil, fmt.Errorf("member must have a model")
	}

	defaultID := uuid.NewString()

	agent := &Agent{
		id:           defaultID,
		name:         defaultID,
		model:        model,
		tools:        make(tools.Toolkit),
		subordinates: make([]AgentImpl, 0),
	}

	for _, opt := range options {
		if err := opt.Apply(agent); err != nil {
			return nil, err
		}
	}

	if len(agent.subordinates) > 0 {
		addDelegateTaskTool(agent)
	}

	return agent, nil
}

func WithID(id string) AgentOption {
	if id == "" {
		id = uuid.NewString()
	}
	return AgentOptionFunc(func(a *Agent) error {
		a.id = id
		return nil
	})
}

func WithName(name string) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		if name != "" {
			a.name = name
		}
		return nil
	})
}

func WithRole(role string) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.role = role
		return nil
	})
}

func WithDescription(description string) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.description = description
		return nil
	})
}

func WithInstructions(instructions string) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.instructions = instructions
		return nil
	})
}

func WithGoal(goal string) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.goal = goal
		return nil
	})
}

func WithExpectedOutput(expectedOutput string) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.expectedOutput = expectedOutput
		return nil
	})
}

func WithTool(tool tools.Tool) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.tools.AddTool(tool)
		return nil
	})
}

func WithRetry(retryOptions *retry.Options) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		if retryOptions != nil {
			a.retry = retryOptions
		}
		return nil
	})
}

func WithSubordinate(subordinate AgentImpl) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.subordinates = append(a.subordinates, subordinate)
		return nil
	})
}

func WithChatOptions(chatOptions provider.ChatOptions) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.chatOptions = chatOptions
		return nil
	})
}
