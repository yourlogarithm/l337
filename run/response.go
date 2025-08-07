package run

import "github.com/yourlogarithm/l337/chat"

type Response struct {
	Messages []chat.Message
}

func (r *Response) AddMessage(msg chat.Message) {
	r.Messages = append(r.Messages, msg)
}

func (r *Response) Content() string {
	if len(r.Messages) == 0 {
		return ""
	}
	return r.Messages[len(r.Messages)-1].Content
}
