package agent

import (
	"context"

	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

type AgentImpl interface {
	Name() (string, error)
	Description() (string, error)
	Skills() ([]tools.SkillCard, error)
	Run(context.Context, *run.Response) error
}
