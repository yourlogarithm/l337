package team

import (
	"fmt"
	"strconv"
	"strings"
)

func (t *Team) computeSystemMessage() (string, error) {
	var sb strings.Builder

	sb.WriteString("You are the leader of a team and sub-teams of AI Agents.\n")
	sb.WriteString("Your task is to coordinate the team to complete the user's request.\n\n")

	sb.WriteString("Here are the members in your team:\n<team_members>\n")
	for i, member := range t.Members {
		options := member.GetOptions()
		sb.WriteString(" - Agent ")
		sb.WriteString(strconv.Itoa(i + 1))
		sb.WriteString(":\n")
		sb.WriteString("   - ID: ")
		sb.WriteString(options.ID)
		sb.WriteString("\n")
		sb.WriteString("   - Name: ")
		sb.WriteString(options.Name)
		sb.WriteString("\n")
		if options.Role != "" {
			sb.WriteString("   - Role: ")
			sb.WriteString(options.Role)
			sb.WriteString("\n")
		}

		if len(options.Tools) > 0 {
			sb.WriteString("   - Member tools:\n")
			for _, tool := range options.Tools {
				sb.WriteString("    - ")
				sb.WriteString(tool.Name)
				sb.WriteString("\n")
			}
		}

		if i < len(t.Members)-1 {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("</team_members>\n\n")

	sb.WriteString("<how_to_respond>\n")
	switch t.Mode {
	case ModeCollaborate:
		sb.WriteString("You can either respond directly or transfer tasks to members in your team with the highest likelihood of completing the user's request.\n")
		sb.WriteString("- Carefully analyze the tools available to the members and their roles before transferring tasks.\n")
		sb.WriteString("- You cannot use a member tool directly. You can only transfer tasks to members.\n")
		sb.WriteString("- When you transfer a task to another member, make sure to include:\n")
		sb.WriteString("  - member_id (str): The ID of the member to transfer the task to. Use only the ID of the member, not the ID of the team followed by the ID of the member.\n")
		sb.WriteString("  - task_description (str): A clear description of the task.\n")
		sb.WriteString("  - expected_output (str): The expected output.\n")
		sb.WriteString("- You can transfer tasks to multiple members at once.\n")
		sb.WriteString("- You must always analyze the responses from members before responding to the user.\n")
		sb.WriteString("- After analyzing the responses from the members, if you feel the task has been completed, you can stop and respond to the user.\n")
		sb.WriteString("- If you are not satisfied with the responses from the members, you should re-assign the task.\n")
	case ModeCoordinate:
		sb.WriteString("You can either respond directly or use the `run_member_agents` tool to run all members in your team to get a collaborative response.\n")
		sb.WriteString("- To run the members in your team, call `run_member_agents` ONLY once. This will run all members in your team.\n")
		sb.WriteString("- Analyze the responses from all members and evaluate whether the task has been completed.\n")
		sb.WriteString("- If you feel the task has been completed, you can stop and respond to the user.\n")
	case ModeRoute:
		sb.WriteString("You can either respond directly or forward tasks to members in your team with the highest likelihood of completing the user's request.\n")
		sb.WriteString("- Carefully analyze the tools available to the members and their roles before forwarding tasks.\n")
		sb.WriteString("- When you forward a task to another Agent, make sure to include:\n")
		sb.WriteString("  - member_id (str): The ID of the member to forward the task to. Use only the ID of the member, not the ID of the team followed by the ID of the member.\n")
		sb.WriteString("  - expected_output (str): The expected output.\n")
		sb.WriteString("- You can forward tasks to multiple members at once.\n")
	default:
		return "", fmt.Errorf("unknown team mode: %s", t.Mode)
	}
	sb.WriteString("</how_to_respond>\n\n")

	sb.WriteString(t.Options.ComputeSystemMessage())

	return sb.String(), nil
}
