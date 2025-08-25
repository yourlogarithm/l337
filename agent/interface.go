package agent

import (
	"context"

	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

type AgentImpl interface {
	Name() string
	Description() string
	Skills() []tools.SkillCard
	run(context.Context, *run.Response) error
}
