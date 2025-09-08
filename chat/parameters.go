package chat

import (
	"github.com/google/uuid"
)

type Parameters struct {
	Messages  []Message
	SessionID uuid.UUID
}

type Parameter interface {
	Apply(*Parameters) error
}

type ParameterFunc func(*Parameters) error

func (s ParameterFunc) Apply(r *Parameters) error { return s(r) }

func WithSessionID(sessionID uuid.UUID) Parameter {
	return ParameterFunc(func(p *Parameters) error {
		p.SessionID = sessionID
		return nil
	})
}

func WithMessage(role Role, content string) Parameter {
	return ParameterFunc(func(p *Parameters) error {
		p.Messages = append(p.Messages, Message{
			Role:    role,
			Content: content,
		})
		return nil
	})
}
