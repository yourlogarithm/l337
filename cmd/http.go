package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
)

type LoggingRoundTripper struct {
	transport http.RoundTripper
	mu        sync.Mutex
	logs      []LoggedRequest
}

type LoggedRequest struct {
	Method   string         `json:"method"`
	URL      string         `json:"url"`
	Body     map[string]any `json:"body"`
	Header   http.Header    `json:"header"`
	RespBody map[string]any `json:"respBody"`
}

func (lrt *LoggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var bodyCopy []byte
	if req.Body != nil {
		bodyCopy, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyCopy))
	}

	parsedBody := make(map[string]any)
	if len(bodyCopy) > 0 {
		if err := json.Unmarshal(bodyCopy, &parsedBody); err != nil {
			parsedBody = map[string]any{"error": "invalid JSON body"}
		}
	}

	resp, err := lrt.transport.RoundTrip(req)
	if err != nil {
		return resp, err
	}

	respBody := make(map[string]any)
	if resp != nil && resp.Body != nil {
		respCopy, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(respCopy, &respBody); err != nil {
			respBody = map[string]any{"error": "invalid JSON response"}
		}
		resp.Body = io.NopCloser(bytes.NewBuffer(respCopy))
	}

	lrt.mu.Lock()
	lrt.logs = append(lrt.logs, LoggedRequest{
		Method:   req.Method,
		URL:      req.URL.String(),
		Body:     parsedBody,
		Header:   req.Header,
		RespBody: respBody,
	})
	lrt.mu.Unlock()

	return resp, nil
}

func (lrt *LoggingRoundTripper) SaveToFile(path string) error {
	lrt.mu.Lock()
	defer lrt.mu.Unlock()
	data, err := json.MarshalIndent(lrt.logs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func newLoggingHTTPClient() (*http.Client, *LoggingRoundTripper) {
	rt := &LoggingRoundTripper{transport: http.DefaultTransport}
	return &http.Client{Transport: rt}, rt
}
