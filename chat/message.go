package chat

type Message struct {
	Role       string `json:"role"`
	Name       string `json:"name"`
	Content    string `json:"content"`
	ToolCallID string `json:"tool_call_id,omitempty"`
}
