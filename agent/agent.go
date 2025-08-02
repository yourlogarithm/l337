package agent

import (
	"context"

	"github.com/yourlogarithm/golagno/chat"
	"github.com/yourlogarithm/golagno/provider"
	"github.com/yourlogarithm/golagno/retry"
	"github.com/yourlogarithm/golagno/run"
	"github.com/yourlogarithm/golagno/tools"
)

type Agent struct {
	Name         string
	Role         string
	Description  string
	Instructions string

	Model *provider.Model

	Tools []tools.Tool

	RetryOptions *retry.Options
}

func (a *Agent) Run(ctx context.Context, messages []chat.Message) (runResponse run.Response, err error) {
	var chatResponse chat.Response
	if a.RetryOptions == nil {
		a.RetryOptions = retry.Default()
	}

	req := chat.Request{
		Messages: messages,
		Tools:    a.Tools,
	}

	if err = a.RetryOptions.Execute(func() error {
		response, err := a.Model.Impl.Chat(ctx, &req)
		if err != nil {
			return err
		}
		chatResponse = response
		return nil
	}); err != nil {
		return runResponse, err
	}
	runResponse.History = append(runResponse.History, chatResponse)

	return runResponse, nil
}
