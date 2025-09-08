package client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func (c *RemoteAgent) RunStreaming(ctx context.Context, runResponse *chat.RunResponse, bufferSize int) (provider.ResponseChannel, error) {
	body, err := json.Marshal(runResponse)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/run_streaming", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	ch := provider.NewResponseChannel(bufferSize)

	go func() {
		defer close(ch)
		defer resp.Body.Close()

		reader := bufio.NewReader(resp.Body)

		parseLine := func(prefix string) (string, bool) {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					ch.SendErr(err)
				}
				return "", false
			}

			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, prefix) {
				return "", true
			}

			payload := strings.TrimPrefix(line, prefix)
			payload = strings.TrimSpace(payload)

			return payload, true
		}

		for {
			event, ok := parseLine("event:")
			if !ok {
				break
			} else if event == "" {
				continue
			}

			data, ok := parseLine("data:")
			if !ok {
				break
			} else if data == "" {
				continue
			}

			switch event {
			case "chunk":
				var responseChunk provider.ResponseChunkMarshalable
				if err := json.Unmarshal([]byte(data), &responseChunk); err != nil {
					ch.SendErr(err)
					return
				}
				if responseChunk.Error != "" {
					ch.SendErr(ErrChunkMessage{Message: responseChunk.Error})
					return
				} else {
					ch.Send(responseChunk.Chunk)
				}
			case "response":
				if err := json.Unmarshal([]byte(data), runResponse); err != nil {
					ch.SendErr(err)
					return
				}
			default:
				ch.SendErr(ErrUnknownEvent{Event: event})
				return
			}
		}
	}()

	return ch, nil
}

func (c *RemoteAgent) RunStreamingWithParams(ctx context.Context, bufferSize int, params ...chat.Parameter) (*chat.RunResponse, provider.ResponseChannel, error) {
	runResponse, err := agent.BuildRunResponse(params...)
	if err != nil {
		return nil, nil, err
	}
	stream, err := c.RunStreaming(ctx, runResponse, bufferSize)
	return runResponse, stream, err
}
