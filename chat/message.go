package chat

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
	// Assistant tool calls
	ToolCalls []ToolCall `json:"tool_calls"`
	// OpenAI: An optional name for the participant. Provides the model information to differentiate between participants of the same role.
	Name string `json:"name"`
	// OpenAI: Assistant refusal message
	Refusal string `json:"refusal"`
	// Anthropic: User boolean indicating whether function call resulted in an error.
	IsErr bool `json:"is_err"`
}

type ToolCall struct {
	// Unique identifier for the tool call.
	ID string `json:"id"`
	// Raw LLM arguments.
	Arguments string `json:"arguments"`
	// Tool name
	Name string `json:"name"`
}
