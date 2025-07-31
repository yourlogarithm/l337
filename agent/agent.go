package agent

import (
	"github.com/yourlogarithm/golagno/provider"
	"github.com/yourlogarithm/golagno/retry"
	"github.com/yourlogarithm/golagno/run"
)

type Agent struct {
	Name         string
	Role         string
	Description  string
	Instructions string

	Model *provider.Model

	RetryOptions *retry.Options
}

func (a *Agent) Run() (run.Response, error) {
	return run.Response{}, nil
}
