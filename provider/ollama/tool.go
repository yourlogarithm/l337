package ollama

import (
	"encoding/json"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/tools"
)

func convertTool(t *tools.Tool) api.Tool {
	var parameters struct {
		Type       string                      `json:"type"`
		Defs       any                         `json:"$defs,omitempty"`
		Items      any                         `json:"items,omitempty"`
		Required   []string                    `json:"required"`
		Properties map[string]api.ToolProperty `json:"properties"`
	}

	if t.Schema != nil {
		parameters.Type = t.Schema.Type
		parameters.Defs = t.Schema.Definitions
		parameters.Items = t.Schema.Items
		parameters.Required = t.Schema.Required
		parameters.Properties = make(map[string]api.ToolProperty, t.Schema.Properties.Len())
		marshaled, _ := json.Marshal(t.Schema.Properties)
		json.Unmarshal(marshaled, &parameters.Properties)
	}

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  parameters,
		},
	}
}
