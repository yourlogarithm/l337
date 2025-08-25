package agent

import (
	"context"

	"github.com/yourlogarithm/l337/run"
)

type AgentImpl interface {
	Name() string
	Description() string
	Skills() []Skill
	run(context.Context, *run.Response) error
}
