package agent

import (
	"context"

	"github.com/google/uuid"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

type AgentOption interface {
	Apply(*Agent) error
}

type AgentOptionFunc func(*Agent) error

func (s AgentOptionFunc) Apply(r *Agent) error { return s(r) }

func New(model *provider.Model, options ...AgentOption) (*Agent, error) {
	defaultID := uuid.New()

	agent := &Agent{
		id:           defaultID,
		name:         defaultID.String(),
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

func WithID(id uuid.UUID) AgentOption {
	if id == uuid.Nil {
		id = uuid.New()
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

func WithMCP(ctx context.Context, session *mcp.ClientSession) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.tools.RegisterMCP(ctx, session)
		return nil
	})
}

func WithSubordinate(subordinate AgentImpl) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.subordinates = append(a.subordinates, subordinate)
		return nil
	})
}

func WithChatOptions(chatOptions chat.Options) AgentOption {
	return AgentOptionFunc(func(a *Agent) error {
		a.chatOptions = chatOptions
		return nil
	})
}
