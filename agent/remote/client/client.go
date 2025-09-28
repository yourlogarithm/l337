package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/yourlogarithm/l337/tools"
)

type RemoteAgent struct {
	BaseURL    string
	HttpClient *http.Client
}

func Default(baseUrl string) *RemoteAgent {
	return &RemoteAgent{
		BaseURL:    baseUrl,
		HttpClient: http.DefaultClient,
	}
}

func checkError(resp *http.Response) error {
	if resp.StatusCode != http.StatusOK {
		var content bytes.Buffer
		if _, err := content.ReadFrom(resp.Body); err != nil {
			return ErrServerResponse{
				StatusCode: resp.StatusCode,
				Err:        err,
			}
		}
		return ErrServerResponse{
			StatusCode: resp.StatusCode,
			Err:        errors.New(content.String()),
		}
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

func (c *RemoteAgent) Tools() (tools []tools.Tool, err error) {
	resp, err := c.HttpClient.Get(c.BaseURL + "/tools")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err := checkError(resp); err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&tools); err != nil {
		return nil, err
	}
	return tools, nil
}
