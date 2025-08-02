package agent

import (
	"context"

	"github.com/yourlogarithm/golagno/chat"
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

func (a *Agent) Run(ctx context.Context, messages []chat.Message) (runResponse run.Response, err error) {
	response, err := a.Model.Impl.Chat(ctx, messages)
	if err != nil {
		return runResponse, err
	}

	runResponse.Content = response.Choices[0].Content

	return runResponse, nil
}
