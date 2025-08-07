package agent

import (
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/internal/logging"
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

func NewFromOptions(options agentic.Options) (*Agent, error) {
	if err := options.Initialize(); err != nil {
		return nil, err
	}
	return &Agent{
		Options: options,
	}, nil
}
