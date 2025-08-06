package agentic

import (
	"strings"

	"github.com/google/uuid"
	"github.com/yourlogarithm/golagno/provider"
	"github.com/yourlogarithm/golagno/retry"
	"github.com/yourlogarithm/golagno/tools"
)

type Options struct {
	ID   string
	Name string

	Role        string
	Description string

	Instructions   string
	Goal           string
	ExpectedOutput string

	Model *provider.Model

	Tools tools.Toolkit

	Retry *retry.Options
}

func (o *Options) SetupID() {
	if o.ID == "" {
		o.ID = uuid.NewString()
	}
}

func (o *Options) ComputeSystemMessage() string {
	var sb strings.Builder

	appendSystemString := func(s, tag string) {
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteRune('\n')
			}
			if tag != "" {
				sb.WriteString("<" + tag + ">\n")
			}
			sb.WriteString(s)
			if tag != "" {
				sb.WriteRune('\n')
				sb.WriteString("</" + tag + ">")
			}
		}
	}

	appendSystemString("Your name is "+o.Name, "")
	appendSystemString(o.Description, "")
	appendSystemString(o.Goal, "goal")
	appendSystemString(o.Instructions, "instructions")
	appendSystemString(o.ExpectedOutput, "expected_output")

	return sb.String()
}
