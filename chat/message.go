package chat

type Message struct {
	Role    Role
	Content string
	// Assistant tool calls
	ToolCalls []ToolCall
	// OpenAI: An optional name for the participant. Provides the model information to differentiate between participants of the same role.
	Name string
	// OpenAI: Assistant refusal message
	Refusal string
	// Anthropic: User boolean indicating whether function call resulted in an error.
	IsErr bool
}

type ToolCall struct {
	// Unique identifier for the tool call.
	ID string
	// Key value pairs of arguments to pass to the tool.
	Arguments map[string]any
	// Tool name
	Name string
}
