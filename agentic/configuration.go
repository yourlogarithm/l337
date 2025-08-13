package agentic

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/retry"
	"github.com/yourlogarithm/l337/tools"
)

type Configuration struct {
	initialized bool

	// Unique identifier for the member
	// If not set, a random UUID (stripped to 8 characters) will be generated
	ID string
	// Appended to the system message
	// If not set, an error will be returned as soon as `Initialize` is called
	Name string
	// Provided to the parent member (if any) to help understand the member's purpose
	Role string
	// Description of the member
	// Appends to the system message
	Description string
	// Appended to the system message
	// Wrapped in <instructions> tags
	Instructions string
	// Appended to the system message
	// Wrapped in <goal> tags
	Goal string
	// Appended to the system message
	// Wrapped in <expected_output> tags
	ExpectedOutput string
	// Model used to send LLM requests
	// If not set, an error will be returned as soon as `Initialize` is called
	*provider.Model
	// Tools for the LLM to use
	Tools tools.Toolkit
	// Retry options for the LLM requests
	// If not set, defaults to `retry.Default()`
	Retry *retry.Options

	provider.ChatOptions
}

// Generates a random ID if not set
//
// Checks that mandatory fields are set
func (o *Configuration) Initialize() error {
	if o.initialized {
		return nil
	}

	if o.ID == "" {
		o.ID = uuid.NewString()[0:8]
	}

	if o.Name == "" {
		return fmt.Errorf("member must have a name")
	}

	if o.Model == nil {
		return fmt.Errorf("member must have a model")
	}

	return nil
}

// Generates the system message for the member
func (o *Configuration) ComputeSystemMessage() string {
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
