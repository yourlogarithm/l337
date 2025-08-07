package chat

type Message struct {
	Role      Role   `json:"role"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	ToolCalls []ToolCall
}
