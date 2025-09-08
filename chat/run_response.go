package chat

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/metrics"
)

type RunResponse struct {
	SessionID uuid.UUID                       `json:"session_id"`
	Messages  []Message                       `json:"messages"`
	Metrics   map[uuid.UUID][]metrics.Metrics `json:"metrics"`
}

// Returns the content of the last message in the response.
func (r *RunResponse) Content() string {
	if len(r.Messages) == 0 {
		return ""
	}
	return r.Messages[len(r.Messages)-1].Content
}
