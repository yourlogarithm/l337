package chat

import (
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

type Response struct {
	ID           string
	Created      int64
	Content      string
	Refusal      string
	ToolCalls    []chat.ToolCall
	FinishReason string
	Metrics      metrics.Metrics
}
