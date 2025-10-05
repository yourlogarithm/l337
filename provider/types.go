package provider

import (
	"time"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/tools"
)

type Request struct {
	Messages []chat.Message
	Tools    []tools.Tool
}

type Response struct {
	ID           string          `json:"id"`
	Created      time.Time       `json:"created"`
	Content      chat.Content    `json:"content"`
	Refusal      string          `json:"refusal"`
	Reasoning    string          `json:"reasoning"`
	ToolCalls    []chat.ToolCall `json:"tool_calls"`
	FinishReason string          `json:"finish_reason"`
	Metrics      metrics.Metrics `json:"metrics"`
}
