package chat

import "github.com/yourlogarithm/l337/chat"

type Response struct {
	ID           string
	Created      int64
	Content      string
	Refusal      string
	ToolCalls    []chat.ToolCall
	FinishReason FinishReason
}

type FinishReason string

const (
	FinishReasonStop         FinishReason = "stop"
	FinishReasonLength       FinishReason = "length"
	FinishReasonToolCalls    FinishReason = "tool_calls"
	FinishReasonFunctionCall FinishReason = "content_filter"
)
