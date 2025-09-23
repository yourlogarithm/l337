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

	if t.Parameters != nil {
		parameters.Type = t.Parameters.Type
		parameters.Defs = t.Parameters.Definitions
		parameters.Items = t.Parameters.Items
		parameters.Required = t.Parameters.Required
		parameters.Properties = make(map[string]api.ToolProperty, t.Parameters.Properties.Len())
		marshaled, _ := json.Marshal(t.Parameters.Properties)
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
