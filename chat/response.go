package chat

type Response struct {
	ID      string
	Created int64
	Choices []Choice
}

type FinishReason string

const (
	FinishReasonStop         FinishReason = "stop"
	FinishReasonLength       FinishReason = "length"
	FinishReasonToolCalls    FinishReason = "tool_calls"
	FinishReasonFunctionCall FinishReason = "content_filter"
)

type Choice struct {
	Content      string
	Refusal      string
	ToolCalls    []ToolCall
	FinishReason FinishReason
}

type ToolCall struct {
	ID        string
	Arguments string
	Name      string
}
