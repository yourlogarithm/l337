package types

type ToolCall struct {
	// Unique identifier for the tool call.
	ID string `json:"id"`
	// Raw LLM arguments.
	Arguments string `json:"arguments"`
	// Tool name
	Name string `json:"name"`
}
