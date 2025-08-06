package agentic

import (
	"context"

	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/run"
)

type Member interface {
	Type() MemberType
	GetOptions() *Options
	Run(ctx context.Context, messages []chat.Message) (runResponse run.Response, err error)
}
