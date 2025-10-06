package chat

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/metrics"
)

type RunResponse struct {
	SessionID uuid.UUID                       `json:"session_id"`
	Messages  []Message                       `json:"messages,omitempty"`
	Metrics   map[uuid.UUID][]metrics.Metrics `json:"metrics,omitempty"`
}

// Returns the content of the last message in the response.
func (r *RunResponse) Content() Content {
	if len(r.Messages) == 0 {
		return Content{}
	}
	return r.Messages[len(r.Messages)-1].Content
}
