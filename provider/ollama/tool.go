package ollama

import (
	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/tools"
)

func convertTool(t *tools.Tool) api.Tool {
	parameters := struct {
		Type       string   "json:\"type\""
		Defs       any      "json:\"$defs,omitempty\""
		Items      any      "json:\"items,omitempty\""
		Required   []string "json:\"required\""
		Properties map[string]struct {
			Type        api.PropertyType "json:\"type\""
			Items       any              "json:\"items,omitempty\""
			Description string           "json:\"description\""
			Enum        []any            "json:\"enum,omitempty\""
		} "json:\"properties\""
	}{
		Type:     "object",
		Required: t.Required,
	}

	parameters.Properties = make(map[string]struct {
		Type        api.PropertyType `json:"type"`
		Items       any              `json:"items,omitempty"`
		Description string           `json:"description"`
		Enum        []any            `json:"enum,omitempty"`
	}, len(t.Parameters))

	for name, schema := range t.Parameters {
		ollamaParam := struct {
			Type        api.PropertyType "json:\"type\""
			Items       any              "json:\"items,omitempty\""
			Description string           "json:\"description\""
			Enum        []any            "json:\"enum,omitempty\""
		}{
			Type:        []string{schema.Type},
			Items:       schema.Items,
			Description: schema.Description,
			Enum:        schema.Enum,
		}
		parameters.Properties[name] = ollamaParam
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
