package agent

import (
	"strconv"
	"strings"
)

// Build the system message for the Agent
func (a *Agent) ComputeSystemMessage() (string, error) {
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

	appendSystemString("Your name is "+a.name, "")
	appendSystemString(a.description, "")
	appendSystemString(a.goal, "goal")
	appendSystemString(a.instructions, "instructions")
	appendSystemString(a.expectedOutput, "expected_output")

	if len(a.subordinates) > 0 {
		var subordinatesSb strings.Builder
		subordinatesSb.WriteString("Here are the members in your team:\n")
		for i, subordinate := range a.subordinates {
			subordinatesSb.WriteString(" - Agent ")
			subordinatesSb.WriteString(strconv.Itoa(i + 1))
			subordinatesSb.WriteString(":\n")

			name, err := subordinate.Name()
			if err != nil {
				return "", err
			}
			subordinatesSb.WriteString("   - Name: ")
			subordinatesSb.WriteString(name)

			desc, err := subordinate.Description()
			if err != nil {
				return "", err
			}
			if desc != "" {
				subordinatesSb.WriteString("   - Description: ")
				subordinatesSb.WriteString(desc)
				subordinatesSb.WriteByte('\n')
			}

			skills, err := subordinate.Skills()
			if err != nil {
				return "", err
			}
			if len(skills) > 0 {
				subordinatesSb.WriteString("   - Member tools:\n")
				for _, skill := range skills {
					subordinatesSb.WriteString("    - ")
					subordinatesSb.WriteString(skill.Name)
					subordinatesSb.WriteByte(':')
					subordinatesSb.WriteString(skill.Description)
					subordinatesSb.WriteByte('\n')
				}
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
