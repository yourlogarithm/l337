package chat

import (
	"time"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

type Response struct {
	ID           string
	Created      time.Time
	Content      string
	Refusal      string
	Reasoning    string
	ToolCalls    []chat.ToolCall
	FinishReason string
	Metrics      metrics.Metrics
}
