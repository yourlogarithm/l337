package agentic

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/run"
)

type Member interface {
	Type() MemberType
	GetOptions() *Configuration
	Run(ctx context.Context, messages []chat.Message) (runResponse run.Response, err error)
}
