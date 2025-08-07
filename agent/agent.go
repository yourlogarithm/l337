package agent

import (
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/logging"
)

var logger = logging.SetupLogger("agent")

type Agent struct {
	agentic.Options
}

func (a *Agent) Type() agentic.MemberType {
	return agentic.MemberTypeAgent
}

func (a *Agent) GetOptions() *agentic.Options {
	return &a.Options
}

func NewFromOptions(options agentic.Options) *Agent {
	options.SetupID()
	return &Agent{
		Options: options,
	}
}
