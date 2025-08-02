package run

import "github.com/yourlogarithm/golagno/chat"

type Response struct {
	Content string
	History []chat.Response
}
