package chat

type Response struct {
	ID           string
	Created      int64
	Content      string
	Refusal      string
	ToolCalls    []ToolCall
	FinishReason FinishReason
}

type FinishReason string

const (
	FinishReasonStop         FinishReason = "stop"
	FinishReasonLength       FinishReason = "length"
	FinishReasonToolCalls    FinishReason = "tool_calls"
	FinishReasonFunctionCall FinishReason = "content_filter"
)

type ToolCall struct {
	ID        string
	Arguments map[string]any
	Name      string
}
