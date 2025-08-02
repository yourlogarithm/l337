package chat

type Response struct {
	ID      string
	Created int64
	Content Content
}

type FinishReason string

const (
	FinishReasonStop         FinishReason = "stop"
	FinishReasonLength       FinishReason = "length"
	FinishReasonToolCalls    FinishReason = "tool_calls"
	FinishReasonFunctionCall FinishReason = "content_filter"
)

type Content struct {
	Text         string
	Refusal      string
	ToolCalls    []ToolCall
	FinishReason FinishReason
}

type ToolCall struct {
	ID        string
	Arguments string
	Name      string
}
