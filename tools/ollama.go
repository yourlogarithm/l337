package tools

import (
	"github.com/invopop/jsonschema"
	"github.com/ollama/ollama/api"
)

func (t *Tool) ToOllamaTool() api.Tool {
	defsMap := make(map[jsonschema.ID]jsonschema.Schema)

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
		Required: make([]string, 0, len(t.Parameters)),
	}

	parameters.Properties = make(map[string]struct {
		Type        api.PropertyType `json:"type"`
		Items       any              `json:"items,omitempty"`
		Description string           `json:"description"`
		Enum        []any            `json:"enum,omitempty"`
	}, len(t.Parameters))

	for name, schema := range t.Parameters {
		parameters.Required = append(parameters.Required, name)
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

	defs := make([]jsonschema.Schema, 0, len(defsMap))
	for _, schema := range defsMap {
		defs = append(defs, schema)
	}
	parameters.Defs = defs

	return api.Tool{
		Type: "function",
		Function: api.ToolFunction{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  parameters,
		},
	}
}
