package chat

type Message struct {
	Role      string `json:"role"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	ToolCalls []ToolCall
}
