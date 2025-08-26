package run

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

type Response struct {
	SessionID uuid.UUID                       `json:"session_id"`
	Messages  []chat.Message                  `json:"messages"`
	Metrics   map[uuid.UUID][]metrics.Metrics `json:"metrics"`
}

// Content of the last message in the response.
func (r *Response) Content() string {
	if len(r.Messages) == 0 {
		return ""
	}
	return r.Messages[len(r.Messages)-1].Content
}
