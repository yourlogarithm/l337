package team

import (
	"context"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/tools"
)

func (t *Team) generateRouteTool() tools.Tool {
	callable := func(ctx context.Context, params tools.Params) (string, error) {
		logger.Debug("route.call", "team", t.Configuration.ID, "params", params)

		memberID, err := tools.GetParameter[string](params, "member_id")
		if err != nil {
			return "", err
		}

		expectedOutput, err := tools.GetParameter[string](params, "expected_output")
		if err != nil {
			return "", err
		}

		member := t.findMemberByID(memberID)
		if member == nil {
			return "", ErrMemberNotFound.Error(memberID)
		}

		messages := []chat.Message{
			{
				Role:    chat.RoleUser,
				Content: expectedOutput,
			},
		}

		response, err := (*member).Run(ctx, messages)

		return response.Content(), nil
	}

	forwardTaskToMember := tools.NewTool(
		"forward_task_to_member",
		"Use this function to forward the request to the selected team member.",
		callable,
	)

	tools.AddParameterFromType[string](&forwardTaskToMember, "member_id", "The ID of the member to transfer the task to. Use only the ID of the member, not the ID of the team followed by the ID of the member.", true)
	tools.AddParameterFromType[string](&forwardTaskToMember, "expected_output", "The expected output from the member.", true)

	return forwardTaskToMember
}
