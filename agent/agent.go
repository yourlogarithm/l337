package agent

import (
	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/internal/logging"
)

var logger = logging.SetupLogger("agent")

type Agent struct {
	agentic.Configuration
}

func (a *Agent) Type() agentic.MemberType {
	return agentic.MemberTypeAgent
}

func (a *Agent) GetOptions() *agentic.Configuration {
	return &a.Configuration
}

func NewFromOptions(options agentic.Configuration) (*Agent, error) {
	if err := options.Initialize(); err != nil {
		return nil, err
	}
	return &Agent{
		Configuration: options,
	}, nil
}
