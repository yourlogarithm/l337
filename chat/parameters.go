package chat

import (
	"github.com/google/uuid"
)

type Parameter interface {
	Apply(*RunResponse) error
}

type ParameterFunc func(*RunResponse) error

func (s ParameterFunc) Apply(r *RunResponse) error { return s(r) }

func WithSessionID(sessionID uuid.UUID) Parameter {
	return ParameterFunc(func(r *RunResponse) error {
		r.SessionID = sessionID
		return nil
	})
}

func WithTextMessage(role Role, content string) Parameter {
	return ParameterFunc(func(r *RunResponse) error {
		r.Messages = append(r.Messages, Message{
			Role:    role,
			Content: NewTextContent(content),
		})
		return nil
	})
}

func WithImageUrlMessage(role Role, imageURL string) Parameter {
	return ParameterFunc(func(r *RunResponse) error {
		r.Messages = append(r.Messages, Message{
			Role:    role,
			Content: NewImageUrlContent(imageURL),
		})
		return nil
	})
}

func WithImageContentMessage(role Role, imageContent []byte, format ImageFormat) Parameter {
	return ParameterFunc(func(r *RunResponse) error {
		r.Messages = append(r.Messages, Message{
			Role:    role,
			Content: NewImageContent(imageContent, format),
		})
		return nil
	})
}

func WithAudioContentMessage(role Role, audioContent []byte, format AudioFormat) Parameter {
	return ParameterFunc(func(r *RunResponse) error {
		r.Messages = append(r.Messages, Message{
			Role:    role,
			Content: NewAudioContent(audioContent, format),
		})
		return nil
	})
}
