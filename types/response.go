package types

import (
	"time"

	"github.com/yourlogarithm/l337/metrics"
)

type Response struct {
	ID           string          `json:"id"`
	Created      time.Time       `json:"created"`
	Content      Content         `json:"content"`
	Refusal      string          `json:"refusal"`
	Reasoning    string          `json:"reasoning"`
	ToolCalls    []ToolCall      `json:"tool_calls"`
	FinishReason string          `json:"finish_reason"`
	Metrics      metrics.Metrics `json:"metrics"`
}
