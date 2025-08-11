package chat

type Message struct {
	Role Role
	// An optional name for the participant. Provides the model information to
	// differentiate between participants of the same role.
	// Not all providers support this.
	Name      string
	Content   string
	ToolCalls []ToolCall
	IsErr     bool
}

type ToolCall struct {
	// Unique identifier for the tool call.
	ID string
	// Key value pairs of arguments to pass to the tool.
	Arguments map[string]any
	// Tool name
	Name string
}
