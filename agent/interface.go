package agent

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/run"
)

type AgentImpl interface {
	Name() string
	Description() string
	Skills() []Skill
	Run(context.Context, []chat.Message) (run.Response, error)
}
