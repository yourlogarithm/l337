package team

import (
	"context"
	"strings"
	"sync"

	"github.com/yourlogarithm/l337/agentic"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/tools"
)

func (t *Team) generateCollaborateTool() tools.Tool {
	callable := func(ctx context.Context, params tools.Params) (string, error) {
		logger.Debug("collaborate.call", "team", t.Options.ID, "params", params)

		task_description, err := tools.GetParameter[string](params, "task_description")
		if err != nil {
			return "", err
		}

		expectedOutput, err := tools.GetParameterOptional[string](params, "expected_output")
		if err != nil {
			return "", err
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
				Role:    chat.RoleUser.String(),
				Content: sb.String(),
			},
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		wg.Add(len(t.Members))

		var memberResponses strings.Builder

		for i := range t.Members {
			go func(member agentic.Member) {
				defer wg.Done()
				opts := member.GetOptions()
				response, err := member.Run(ctx, messages)
				var content string
				if err != nil {
					logger.Error("Error running member", "member_id", opts.ID, "error", err)
					content = "Error: " + err.Error()
				} else {
					content = response.Content()
					logger.Debug("Member response", "member_id", opts.ID, "response", content)
				}
				mu.Lock()
				defer mu.Unlock()
				if memberResponses.Len() > 0 {
					memberResponses.WriteString("\n\n")
				}
				memberResponses.WriteString("Member - ")
				memberResponses.WriteString(opts.ID)
				memberResponses.WriteString(":\n")
				memberResponses.WriteString(content)
			}(t.Members[i])
		}

		wg.Wait()

		return memberResponses.String(), nil
	}

	forwardTaskToMember := tools.NewTool(
		"run_member_agents",
		"Send the same task to all the member agents and return the responses.",
		callable,
	)

	tools.AddParameterFromType[string](&forwardTaskToMember, "task_description", "The task description to send to the members.", true)
	tools.AddParameterFromType[string](&forwardTaskToMember, "expected_output", "The expected output from the members (optional).", false)

	return forwardTaskToMember
}
