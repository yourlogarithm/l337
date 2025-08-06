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
	Method string         `json:"method"`
	URL    string         `json:"url"`
	Body   map[string]any `json:"body"`
	Header http.Header    `json:"header"`
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

	lrt.mu.Lock()
	lrt.logs = append(lrt.logs, LoggedRequest{
		Method: req.Method,
		URL:    req.URL.String(),
		Body:   parsedBody,
		Header: req.Header,
	})
	lrt.mu.Unlock()

	return lrt.transport.RoundTrip(req)
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
