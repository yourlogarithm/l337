package remote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/run"
	"github.com/yourlogarithm/l337/tools"
)

type RemoteAgent struct {
	BaseURL    string
	HttpClient *http.Client
}

func DefaultClient(baseUrl string) *RemoteAgent {
	return &RemoteAgent{
		BaseURL:    baseUrl,
		HttpClient: http.DefaultClient,
	}
}

func checkError(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		var content bytes.Buffer
		if _, err := content.ReadFrom(resp.Body); err != nil {
			return fmt.Errorf("unexpected status code - %d: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("unexpected status code - %d: %s", resp.StatusCode, content.String())
	}
	return nil
}

func (c *RemoteAgent) getString(endpoint string) (string, error) {
	resp, err := c.HttpClient.Get(c.BaseURL + endpoint)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *RemoteAgent) Name() (string, error) {
	name, err := c.getString("/name")
	if err != nil {
		return "", err
	}
	return name, nil
}

func (c *RemoteAgent) Description() (string, error) {
	description, err := c.getString("/description")
	if err != nil {
		return "", err
	}
	return description, nil
}

func (c *RemoteAgent) Skills() (skills []tools.SkillCard, err error) {
	resp, err := c.HttpClient.Get(c.BaseURL + "/skills")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&skills); err != nil {
		return nil, err
	}
	return skills, nil
}

func (c *RemoteAgent) Run(ctx context.Context, runResponse *run.Response) error {
	body, err := json.Marshal(runResponse)
	if err != nil {
		return err
	}

	resp, err := c.HttpClient.Post(c.BaseURL+"/run", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return err
	}

	if err := json.NewDecoder(resp.Body).Decode(&runResponse); err != nil {
		return fmt.Errorf("failed to decode run response: %w", err)
	}

	return nil
}

func (r *RemoteAgent) RunWithParams(ctx context.Context, params ...run.Parameter) (run.Response, error) {
	var runParams run.Parameters
	for _, param := range params {
		if err := param.Apply(&runParams); err != nil {
			return run.Response{}, err
		}
	}
	if len(runParams.Messages) == 0 {
		return run.Response{}, fmt.Errorf("no messages provided")
	}

	runResponse := &run.Response{
		SessionID: runParams.SessionID,
		Messages:  runParams.Messages,
		Metrics:   make(map[uuid.UUID][]metrics.Metrics),
	}

	return *runResponse, r.Run(ctx, runResponse)
}
