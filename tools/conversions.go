package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	google "github.com/google/jsonschema-go/jsonschema"
	invopop "github.com/invopop/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"github.com/yourlogarithm/l337/chat"
)

func convertMCPTool(session *mcp.ClientSession, mcpTool *mcp.Tool) Tool {
	schema := googleSchemaToInvopopSchema(mcpTool.InputSchema)
	return Tool{
		Name:        mcpTool.Name,
		Description: mcpTool.Description,
		Parameters:  schema,
		callable: func(ctx context.Context, runResponse *chat.RunResponse, rawArguments string) (string, error) {
			args := map[string]any{}
			if err := json.Unmarshal([]byte(rawArguments), &args); err != nil {
				return "", fmt.Errorf("invalid tool arguments: %w", err)
			}
			params := &mcp.CallToolParams{
				Name:      mcpTool.Name,
				Arguments: args,
			}
			callResult, err := session.CallTool(ctx, params)
			if err != nil {
				return "", err
			} else if callResult.IsError {
				return "", fmt.Errorf("tool call error: %s", callResult.Content)
			}
			var sb strings.Builder
			for i, content := range callResult.Content {
				if textContent, ok := content.(*mcp.TextContent); ok {
					sb.WriteString(textContent.Text)
					if i > 0 {
						sb.WriteByte('\n')
					}
				} else {
					return "", fmt.Errorf("unsupported content type: %T", content)
				}
			}
			return sb.String(), nil
		},
	}
}

func googleSchemaToInvopopSchema(src *google.Schema) *invopop.Schema {
	if src == nil {
		return nil
	}

	dst := &invopop.Schema{
		Version:          src.Schema,
		ID:               invopop.ID(src.ID),
		Anchor:           src.Anchor,
		Ref:              src.Ref,
		DynamicRef:       src.DynamicRef,
		Comments:         src.Comment,
		Title:            src.Title,
		Description:      src.Description,
		Deprecated:       src.Deprecated,
		ReadOnly:         src.ReadOnly,
		WriteOnly:        src.WriteOnly,
		Format:           src.Format,
		Examples:         src.Examples,
		Enum:             src.Enum,
		Const:            derefAny(src.Const),
		ContentEncoding:  src.ContentEncoding,
		ContentMediaType: src.ContentMediaType,
	}

	// Default: google uses json.RawMessage, invopop uses any
	if len(src.Default) > 0 {
		var def any
		_ = json.Unmarshal(src.Default, &def)
		dst.Default = def
	}

	// Types: invopop only supports single Type
	if src.Type != "" {
		dst.Type = src.Type
	} else if len(src.Types) > 0 {
		dst.Type = src.Types[0]
	}

	// Numbers: wrap *float64 into json.Number
	num := func(f *float64) json.Number {
		if f == nil {
			return ""
		}
		return json.Number(fmt.Sprintf("%g", *f))
	}
	dst.MultipleOf = num(src.MultipleOf)
	dst.Maximum = num(src.Maximum)
	dst.Minimum = num(src.Minimum)
	dst.ExclusiveMaximum = num(src.ExclusiveMaximum)
	dst.ExclusiveMinimum = num(src.ExclusiveMinimum)

	// Integers: google uses *int, invopop uses *uint64
	toUint := func(i *int) *uint64 {
		if i == nil {
			return nil
		}
		u := uint64(*i)
		return &u
	}
	dst.MinLength = toUint(src.MinLength)
	dst.MaxLength = toUint(src.MaxLength)
	dst.MinItems = toUint(src.MinItems)
	dst.MaxItems = toUint(src.MaxItems)
	dst.MinProperties = toUint(src.MinProperties)
	dst.MaxProperties = toUint(src.MaxProperties)
	dst.MinContains = toUint(src.MinContains)
	dst.MaxContains = toUint(src.MaxContains)

	dst.Pattern = src.Pattern
	dst.UniqueItems = src.UniqueItems

	// Subschemas
	dst.Items = googleSchemaToInvopopSchema(src.Items)
	dst.Contains = googleSchemaToInvopopSchema(src.Contains)
	dst.PrefixItems = convertSlice(src.PrefixItems)

	dst.AllOf = convertSlice(src.AllOf)
	dst.AnyOf = convertSlice(src.AnyOf)
	dst.OneOf = convertSlice(src.OneOf)
	dst.Not = googleSchemaToInvopopSchema(src.Not)

	dst.If = googleSchemaToInvopopSchema(src.If)
	dst.Then = googleSchemaToInvopopSchema(src.Then)
	dst.Else = googleSchemaToInvopopSchema(src.Else)

	dst.DependentSchemas = convertMap(src.DependentSchemas)
	dst.DependentRequired = src.DependentRequired

	// Object schemas
	if src.Properties != nil {
		props := orderedmap.New[string, *invopop.Schema]()
		for k, v := range src.Properties {
			props.Set(k, googleSchemaToInvopopSchema(v))
		}
		dst.Properties = props
	}
	dst.PatternProperties = convertMap(src.PatternProperties)
	dst.AdditionalProperties = googleSchemaToInvopopSchema(src.AdditionalProperties)
	dst.PropertyNames = googleSchemaToInvopopSchema(src.PropertyNames)

	// Definitions
	if src.Defs != nil {
		dst.Definitions = convertMap(src.Defs)
	} else if src.Definitions != nil {
		dst.Definitions = convertMap(src.Definitions)
	}

	// Extras
	dst.Extras = src.Extra

	return dst
}

func convertSlice(src []*google.Schema) []*invopop.Schema {
	if src == nil {
		return nil
	}
	out := make([]*invopop.Schema, len(src))
	for i, s := range src {
		out[i] = googleSchemaToInvopopSchema(s)
	}
	return out
}

func convertMap(src map[string]*google.Schema) map[string]*invopop.Schema {
	if src == nil {
		return nil
	}
	out := make(map[string]*invopop.Schema, len(src))
	for k, v := range src {
		out[k] = googleSchemaToInvopopSchema(v)
	}
	return out
}

func derefAny(p *any) any {
	if p == nil {
		return nil
	}
	return *p
}
