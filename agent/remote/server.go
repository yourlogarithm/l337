package remote

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/run"
)

type AgentServer struct {
	Agent agent.AgentImpl
}

func returnPlaintext(w http.ResponseWriter, content string) {
	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(content)); err != nil {
		http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
	}
}

func (r *AgentServer) Name(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	name, err := r.Agent.Name()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get agent name: %v", err), http.StatusInternalServerError)
		return
	}

	returnPlaintext(w, name)
}

func (r *AgentServer) Description(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	description, err := r.Agent.Description()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get agent description: %v", err), http.StatusInternalServerError)
		return
	}

	returnPlaintext(w, description)
}

func (r *AgentServer) Skills(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	skills, err := r.Agent.Skills()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get agent skills: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(skills); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func (r *AgentServer) Run(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var runResponse run.Response
	if err := json.NewDecoder(req.Body).Decode(&runResponse); err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	err := r.Agent.Run(ctx, &runResponse)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute run: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(runResponse); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}

func (r *AgentServer) Serve(addr string, handler http.Handler) error {
	http.HandleFunc("/name", r.Name)
	http.HandleFunc("/skills", r.Skills)
	http.HandleFunc("/description", r.Description)
	http.HandleFunc("/run", r.Run)
	return http.ListenAndServe(addr, handler)
}
