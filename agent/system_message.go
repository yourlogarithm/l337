package agent

import (
	"strings"
)

// Build the system message for the Agent
func (a *Agent) ComputeSystemMessage() (string, error) {
	if (a.name == "" || a.name == a.id.String()) && a.description == "" && a.instructions == "" && a.expectedOutput == "" && len(a.subordinates) == 0 {
		// Avoid sending a redundant system message because of the auto-generated UUID.
		return "", nil
	}

	var sb strings.Builder

	appendSystemString := func(s, tag string) {
		if s == "" {
			return
		}
		if tag != "" {
			sb.WriteString("\n\n")
			sb.WriteString("<" + tag + ">\n")
			if strings.Contains(s, "\n") {
				sb.WriteString(s)
			} else {
				sb.WriteString("  " + s)
			}
			sb.WriteRune('\n')
			sb.WriteString("</" + tag + ">")
		} else if sb.Len() > 0 {
			sb.WriteRune('\n')
			sb.WriteString(s)
		} else {
			sb.WriteString(s)
		}
	}

	if a.name != "" {
		sb.WriteString("Your name is " + a.name + ".")
	}
	appendSystemString(a.description, "")
	appendSystemString(a.instructions, "instructions")
	appendSystemString(a.expectedOutput, "expected_output")

	if len(a.subordinates) > 0 {
		var subordinatesSb strings.Builder
		for i, subordinate := range a.subordinates {
			subordinatesSb.WriteString("  - Name: ")
			name, err := subordinate.Name()
			if err != nil {
				return "", err
			}
			subordinatesSb.WriteString(name)
			subordinatesSb.WriteString("\n")

			desc, err := subordinate.Description()
			if err != nil {
				return "", err
			}
			if desc != "" {
				subordinatesSb.WriteString("    Description: ")
				subordinatesSb.WriteString(desc)
			}

			if i < len(a.subordinates)-1 {
				subordinatesSb.WriteByte('\n')
			}
		}
		appendSystemString(subordinatesSb.String(), "subordinates")

		var taskDelegationSb strings.Builder
		taskDelegationSb.WriteString("Depending on the nature of the user request, you can choose to delegate tasks to one or more of your subordinates and then synthesize their responses, or respond directly to the user.")
		appendSystemString(taskDelegationSb.String(), "task_delegation")
	}

	return sb.String(), nil
}
