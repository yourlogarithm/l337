package chat

import (
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/tools"
)

type Request struct {
	Messages []chat.Message
	Tools    []tools.Tool
}
