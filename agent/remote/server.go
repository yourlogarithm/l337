package remote

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/internal/logging"
	"github.com/yourlogarithm/l337/run"
)

var serverLogger = logging.SetupLogger("remote.server")

type AgentServer struct {
	Agent agent.AgentImpl
}

func returnPlaintext(w http.ResponseWriter, req *http.Request, content string) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(content)); err != nil {
		serverLogger.Error("Failed to write response", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
	}
}

func (r *AgentServer) Name(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Initiated", "method", req.Method, "path", req.URL.Path)

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

	serverLogger.Info("Completed", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Description(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Initiated", "method", req.Method, "path", req.URL.Path)

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

	serverLogger.Info("Completed", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Skills(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Initiated", "method", req.Method, "path", req.URL.Path)

	if req.Method != http.MethodGet {
		serverLogger.Info("Invalid request method", "method", req.Method, "path", req.URL.Path)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	skills, err := r.Agent.Skills()
	if err != nil {
		serverLogger.Error("Failed to get agent skills", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to get agent skills: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(skills); err != nil {
		serverLogger.Error("Failed to encode response", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}

	serverLogger.Info("Completed", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Run(w http.ResponseWriter, req *http.Request) {
	serverLogger.Debug("Initiated", "method", req.Method, "path", req.URL.Path)

	if req.Method != http.MethodPost {
		serverLogger.Info("Invalid request method", "method", req.Method, "path", req.URL.Path)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var runResponse run.Response
	if err := json.NewDecoder(req.Body).Decode(&runResponse); err != nil {
		serverLogger.Info("Failed to parse request body", "error", err, "path", req.URL.Path)
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	err := r.Agent.Run(ctx, &runResponse)
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

	serverLogger.Info("Completed", "method", req.Method, "path", req.URL.Path)
}

func (r *AgentServer) Serve(addr string, handler http.Handler) error {
	http.HandleFunc("/name", r.Name)
	http.HandleFunc("/skills", r.Skills)
	http.HandleFunc("/description", r.Description)
	http.HandleFunc("/run", r.Run)
	return http.ListenAndServe(addr, handler)
}
