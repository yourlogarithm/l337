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

func WithTextMessage(role Role, content string) Parameter {
	return ParameterFunc(func(p *Parameters) error {
		p.Messages = append(p.Messages, Message{
			Role:    role,
			Content: NewTextContent(content),
		})
		return nil
	})
}

func WithImageUrlMessage(role Role, imageURL string) Parameter {
	return ParameterFunc(func(p *Parameters) error {
		p.Messages = append(p.Messages, Message{
			Role:    role,
			Content: NewImageUrlContent(imageURL),
		})
		return nil
	})
}

func WithImageContentMessage(role Role, imageContent []byte) Parameter {
	return ParameterFunc(func(p *Parameters) error {
		p.Messages = append(p.Messages, Message{
			Role:    role,
			Content: NewImageContent(imageContent),
		})
		return nil
	})
}
