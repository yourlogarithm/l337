package chat

import "github.com/yourlogarithm/l337/tools"

type Request struct {
	Messages []Message
	Tools    []tools.Tool
}
