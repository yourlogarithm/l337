package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/tools"
)

func TestAddDelegateTaskTool(t *testing.T) {
	t.Run("successful delegation tool addition", func(t *testing.T) {
		subordinate := &mockAgent{
			name:        "test-subordinate",
			description: "A test subordinate",
			skills:      []tools.SkillCard{},
		}

		agent, err := New(
			createMockModel(),
			WithName("manager"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)

		// Check that delegate_task tool was added automatically
		skills, err := agent.Skills()
		require.NoError(t, err)

		assert.Equal(t, 1, len(skills)) // Only the delegate_task tool should be present

		assert.Equal(t, "delegate_task", skills[0].Name)
		assert.Equal(t, "Delegates the task to one or more subordinates", skills[0].Description)
	})

	t.Run("no delegate tool without subordinates", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("solo-agent"))
		require.NoError(t, err)

		skills, err := agent.Skills()
		require.NoError(t, err)

		// Should not have delegate_task tool without subordinates
		assert.Equal(t, 0, len(skills))
	})

	t.Run("error with nil agent", func(t *testing.T) {
		err := addDelegateTaskTool(nil)
		assert.Error(t, err)

		var builderErr ErrBuilderParams
		assert.ErrorAs(t, err, &builderErr)
		assert.Equal(t, "agent", builderErr.Param)
		assert.Equal(t, "agent is required", builderErr.Msg)
	})
}
