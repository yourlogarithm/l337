package tools

import (
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/ollama/ollama/api"
)

func (t *Tool) ToOllamaTool() (api.Tool, error) {
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

	var updateDefs func(e *Element)
	updateDefs = func(e *Element) {
		if e.Schema != nil {
			defsMap[e.Schema.ID] = *e.Schema
		}
		if e.Nested != nil {
			updateDefs(e.Nested)
		}
	}

	for _, param := range t.Parameters {
		parameters.Required = append(parameters.Required, param.Name)
		ollamaParam := struct {
			Type        api.PropertyType "json:\"type\""
			Items       any              "json:\"items,omitempty\""
			Description string           "json:\"description\""
			Enum        []any            "json:\"enum,omitempty\""
		}{
			Type: []string{param.Element.Type},
		}
		if param.Element.Nested != nil {
			ollamaParam.Items = param.Element.Nested.Type
			if param.Element.Nested.Nested != nil {
				return api.Tool{}, fmt.Errorf("ToOllamaTool: Ollama api does not support double nested 2D arrays or maps")
			}
		}
		updateDefs(&param.Element)
		parameters.Properties[param.Name] = ollamaParam
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
	}, nil
}
