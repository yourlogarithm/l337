package google

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/tools"
	"google.golang.org/genai"
)

func uint64ToInt64(u *uint64) *int64 {
	if u == nil {
		return nil
	}
	i := int64(*u)
	return &i
}

func convertSchema(schema *jsonschema.Schema) (genai.Schema, error) {
	out := genai.Schema{
		Default:       schema.Default,
		Description:   schema.Description,
		Example:       schema.Examples,
		Format:        schema.Format,
		MaxItems:      uint64ToInt64(schema.MaxItems),
		MaxLength:     uint64ToInt64(schema.MaxLength),
		MaxProperties: uint64ToInt64(schema.MaxProperties),
		MinItems:      uint64ToInt64(schema.MinItems),
		MinLength:     uint64ToInt64(schema.MinLength),
		MinProperties: uint64ToInt64(schema.MinProperties),
		Pattern:       schema.Pattern,
		Required:      schema.Required,
		Title:         schema.Title,
		Type:          genai.Type(strings.ToUpper(schema.Type)),
	}
	for _, anyOf := range schema.AnyOf {
		converted, err := convertSchema(anyOf)
		if err != nil {
			return genai.Schema{}, err
		}
		out.AnyOf = append(out.AnyOf, &converted)
	}
	for _, enum := range schema.Enum {
		if v, ok := enum.(string); ok {
			out.Enum = append(out.Enum, v)
		} else {
			out.Enum = append(out.Enum, fmt.Sprintf("%v", enum))
		}
	}
	if schema.Items != nil {
		converted, err := convertSchema(schema.Items)
		if err != nil {
			return genai.Schema{}, err
		}
		out.Items = &converted
	}
	if schema.Maximum != "" {
		converted, err := schema.Maximum.Float64()
		if err != nil {
			return genai.Schema{}, err
		}
		out.Maximum = &converted
	}
	if schema.Minimum != "" {
		converted, err := schema.Minimum.Float64()
		if err != nil {
			return genai.Schema{}, err
		}
		out.Minimum = &converted
	}
	if schema.Properties.Len() > 0 {
		bytes, err := schema.Properties.MarshalJSON()
		if err != nil {
			return genai.Schema{}, err
		}
		if err = json.Unmarshal(bytes, &out.Properties); err != nil {
			return genai.Schema{}, err
		}
	}

	return out, nil
}

func convertTool(t *tools.Tool) (tool genai.Tool, err error) {
	schema, err := convertSchema(t.Schema)
	if err != nil {
		return genai.Tool{}, err
	}
	tool = genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:                 t.Name,
				Description:          t.Description,
				ParametersJsonSchema: &schema,
				Response: &genai.Schema{
					Type: genai.TypeString,
				},
			},
		},
	}
	return
}

func convertToolCall(tc *chat.ToolCall) (genai.FunctionCall, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(tc.Arguments), &args); err != nil {
		return genai.FunctionCall{}, err
	}

	return genai.FunctionCall{
		ID:   tc.ID,
		Name: tc.Name,
		Args: args,
	}, nil
}
