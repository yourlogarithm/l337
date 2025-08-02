package chat

import "github.com/yourlogarithm/golagno/tools"

type Request struct {
	Messages []Message
	Tools    []tools.Tool
}
