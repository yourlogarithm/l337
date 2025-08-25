package run

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

type Response struct {
	SessionID uuid.UUID
	Messages  []chat.Message
	Metrics   map[uuid.UUID][]metrics.Metrics
}

// Content of the last message in the response.
func (r *Response) Content() string {
	if len(r.Messages) == 0 {
		return ""
	}
	return r.Messages[len(r.Messages)-1].Content
}
