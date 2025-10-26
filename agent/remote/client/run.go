package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/types"
)

func (c *RemoteAgent) Run(ctx context.Context, run *types.Run) error {
	body, err := json.Marshal(run)
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

	if err := json.NewDecoder(resp.Body).Decode(&run); err != nil {
		return ErrServerResponse{
			Err:        err,
			StatusCode: resp.StatusCode,
		}
	}

	return nil
}

func (c *RemoteAgent) RunWithParams(ctx context.Context, params ...types.Parameter) (types.Run, error) {
	run := &types.Run{
		Metrics: make(map[uuid.UUID][]metrics.Metrics),
	}
	for _, param := range params {
		if err := param.Apply(run); err != nil {
			return types.Run{}, err
		}
	}
	if len(run.Messages) == 0 {
		return types.Run{}, agent.ErrBuilderParams{Param: "Messages", Msg: "at least one message is required to run chat"}
	}
	return *run, c.Run(ctx, run)
}
