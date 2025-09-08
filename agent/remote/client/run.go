package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

func (c *RemoteAgent) Run(ctx context.Context, runResponse *chat.RunResponse) error {
	body, err := json.Marshal(runResponse)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/run", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&runResponse); err != nil {
		return ErrServerResponse{
			Err:        err,
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *RemoteAgent) RunWithParams(ctx context.Context, params ...chat.Parameter) (chat.RunResponse, error) {
	var runParams chat.Parameters
	for _, param := range params {
		if err := param.Apply(&runParams); err != nil {
			return chat.RunResponse{}, err
		}
	}
	if len(runParams.Messages) == 0 {
		return chat.RunResponse{}, agent.ErrBuilderParams{Param: "Messages", Msg: "at least one message is required to run chat"}
	}

	runResponse := &chat.RunResponse{
		SessionID: runParams.SessionID,
		Messages:  runParams.Messages,
		Metrics:   make(map[uuid.UUID][]metrics.Metrics),
	}

	return *runResponse, c.Run(ctx, runResponse)
}
