package team

import (
	"context"
	"strings"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/tools"
)

func (t *Team) generateCoordinateTool() tools.Tool {
	callable := func(ctx context.Context, params tools.Params) (string, error) {
		logger.Debug("coordinate.call", "team", t.Configuration.ID, "params", params)

		memberID, err := tools.GetParameter[string](params, "member_id")
		if err != nil {
			return "", err
		}

		task_description, err := tools.GetParameter[string](params, "task_description")
		if err != nil {
			return "", err
		}

		expectedOutput, err := tools.GetParameterOptional[string](params, "expected_output")
		if err != nil {
			return "", err
		}

		member := t.findMemberByID(memberID)
		if member == nil {
			return "", ErrMemberNotFound.Error(memberID)
		}

		var sb strings.Builder

		sb.WriteString(task_description)
		if expectedOutput != "" {
			sb.WriteString("\n<expected_output>\n")
			sb.WriteString(expectedOutput)
			sb.WriteString("\n</expected_output>")
		}

		messages := []chat.Message{
			{
				Role:    chat.RoleUser,
				Content: sb.String(),
			},
		}

		response, err := (*member).Run(ctx, messages)

		return response.Content(), nil
	}

	forwardTaskToMember := tools.NewTool(
		"transfer_task_to_member",
		"Use this function to transfer a task to the selected team member.",
		callable,
	)

	tools.AddParameterFromType[string](&forwardTaskToMember, "member_id", "The ID of the member to transfer the task to. Use only the ID of the member, not the ID of the team followed by the ID of the member.", true)
	tools.AddParameterFromType[string](&forwardTaskToMember, "task_description", "A clear and concise description of the task the member should achieve.", true)
	tools.AddParameterFromType[string](&forwardTaskToMember, "expected_output", "The expected output from the member (optional).", false)

	return forwardTaskToMember
}
