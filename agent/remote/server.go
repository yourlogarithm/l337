package remote

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/internal/logging"
)

var serverLogger = logging.SetupLogger("remote.server")

type AgentServer struct {
	Agent               agent.AgentImpl
	StreamingBufferSize int
}

func returnPlaintext(w http.ResponseWriter, req *http.Request, content string) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(content)); err != nil {
		serverLogger.Error("Failed to write response", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
	}
}

func (r *AgentServer) Name(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Name", "method", req.Method, "path", req.URL.Path)

	if req.Method != http.MethodGet {
		serverLogger.Info("Invalid request method", "method", req.Method, "path", req.URL.Path)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name, err := r.Agent.Name()
	if err != nil {
		serverLogger.Error("Failed to get agent name", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to get agent name: %v", err), http.StatusInternalServerError)
		return
	}

	returnPlaintext(w, req, name)

	serverLogger.Info("Name", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Description(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Description", "method", req.Method, "path", req.URL.Path)

	if req.Method != http.MethodGet {
		serverLogger.Info("Invalid request method", "method", req.Method, "path", req.URL.Path)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	description, err := r.Agent.Description()
	if err != nil {
		serverLogger.Error("Failed to get agent description", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to get agent description: %v", err), http.StatusInternalServerError)
		return
	}

	returnPlaintext(w, req, description)

	serverLogger.Info("Description", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Tools(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Tools", "method", req.Method, "path", req.URL.Path)

	if req.Method != http.MethodGet {
		serverLogger.Info("Invalid request method", "method", req.Method, "path", req.URL.Path)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	tools, err := r.Agent.Tools()
	if err != nil {
		serverLogger.Error("Failed to get agent tools", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to get agent tools: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tools); err != nil {
		serverLogger.Error("Failed to encode response", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}

	serverLogger.Info("Tools", "method", req.Method, "path", req.URL.Path)
}

func parseRunRequest(w http.ResponseWriter, req *http.Request) *chat.RunResponse {
	if req.Method != http.MethodPost {
		serverLogger.Info("Invalid request method", "method", req.Method, "path", req.URL.Path)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return nil
	}

	var runResponse chat.RunResponse
	if err := json.NewDecoder(req.Body).Decode(&runResponse); err != nil {
		serverLogger.Info("Failed to parse request body", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return nil
	}

	return &runResponse
}

func (r *AgentServer) Run(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Run", "method", req.Method, "path", req.URL.Path)

	runResponse := parseRunRequest(w, req)
	if runResponse == nil {
		return
	}

	ctx := req.Context()
	err := r.Agent.Run(ctx, runResponse)
	if err != nil {
		serverLogger.Error("Failed to execute run", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to execute run: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(runResponse); err != nil {
		serverLogger.Error("Failed to encode response", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}

	serverLogger.Info("Run", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) RunStreaming(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("RunStreaming", "method", req.Method, "path", req.URL.Path)

	runResponse := parseRunRequest(w, req)
	if runResponse == nil {
		return
	}

	ctx := req.Context()
	stream, err := r.Agent.RunStreaming(ctx, runResponse, r.StreamingBufferSize)
	if err != nil {
		serverLogger.Error("Failed to execute run streaming", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to execute run streaming: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	clientGone := req.Context().Done()

	rc := http.NewResponseController(w)

	sendEvent := func(event string, structuredData any) bool {
		data, err := json.Marshal(structuredData)
		if err != nil {
			serverLogger.Error("Failed to marshal data", "error", err, "path", req.URL.Path)
			http.Error(w, fmt.Sprintf("Failed to marshal data: %v", err), http.StatusInternalServerError)
			return false
		}
		_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, data)
		if err != nil {
			serverLogger.Info("Failed to write data", "error", err, "path", req.URL.Path)
			return false
		}
		err = rc.Flush()
		if err != nil {
			serverLogger.Debug("Client disconnected during flush", "error", err, "path", req.URL.Path)
			return false
		}
		return true
	}

	ok := true
streamingLoop:
	for {
		select {
		case <-clientGone:
			serverLogger.Debug("Client disconnected", "path", req.URL.Path)
			break streamingLoop
		case responseChunk, ok := <-stream:
			if !ok {
				serverLogger.Debug("Stream closed by server", "path", req.URL.Path)
				break streamingLoop
			}
			if ok = sendEvent("chunk", responseChunk.ToMarshalable()); !ok {
				break streamingLoop
			}
		}
	}

	if ok {
		sendEvent("response", runResponse)
	}

	serverLogger.Info("RunStreaming", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Serve(addr string, handler http.Handler) error {
	http.HandleFunc("/name", r.Name)
	http.HandleFunc("/tools", r.Tools)
	http.HandleFunc("/description", r.Description)
	http.HandleFunc("/run", r.Run)
	http.HandleFunc("/run_streaming", r.RunStreaming)
	return http.ListenAndServe(addr, handler)
}
