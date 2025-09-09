package agent

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

// mockAgent implements AgentImpl for testing subordinate agents
type mockAgent struct {
	name        string
	description string
	skills      []tools.SkillCard
	nameErr     error
	descErr     error
	skillsErr   error
}

func (m *mockAgent) Name() (string, error) {
	return m.name, m.nameErr
}

func (m *mockAgent) Description() (string, error) {
	return m.description, m.descErr
}

func (m *mockAgent) Skills() ([]tools.SkillCard, error) {
	return m.skills, m.skillsErr
}

func (m *mockAgent) Run(ctx context.Context, runResponse *chat.RunResponse) error {
	return nil
}

func (m *mockAgent) RunStreaming(ctx context.Context, runResponse *chat.RunResponse, bufferSize int) (provider.ResponseChannel, error) {
	return nil, nil
}

func TestComputeSystemMessage(t *testing.T) {
	t.Run("basic system message with name only", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("test-agent"))
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		expected := "Your name is test-agent"
		assert.Equal(t, expected, msg)
	})
	
	t.Run("system message with name and description", func(t *testing.T) {
		agent, err := New(
			createMockModel(),
			WithName("test-agent"),
			WithDescription("A test agent for unit testing"),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		assert.Contains(t, msg, "Your name is test-agent")
		assert.Contains(t, msg, "A test agent for unit testing")
	})
	
	t.Run("system message with all fields", func(t *testing.T) {
		agent, err := New(
			createMockModel(),
			WithName("comprehensive-agent"),
			WithDescription("A comprehensive test agent"),
			WithGoal("Complete all test tasks"),
			WithInstructions("Follow all test protocols"),
			WithExpectedOutput("Detailed test results"),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		// Check all components are present
		assert.Contains(t, msg, "Your name is comprehensive-agent")
		assert.Contains(t, msg, "A comprehensive test agent")
		assert.Contains(t, msg, "<goal>\nComplete all test tasks\n</goal>")
		assert.Contains(t, msg, "<instructions>\nFollow all test protocols\n</instructions>")
		assert.Contains(t, msg, "<expected_output>\nDetailed test results\n</expected_output>")
	})
	
	t.Run("system message with empty fields", func(t *testing.T) {
		agent, err := New(
			createMockModel(),
			WithName("minimal-agent"),
			WithDescription(""),
			WithGoal(""),
			WithInstructions(""),
			WithExpectedOutput(""),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		// Only name should be present
		assert.Equal(t, "Your name is minimal-agent", msg)
		assert.NotContains(t, msg, "<goal>")
		assert.NotContains(t, msg, "<instructions>")
		assert.NotContains(t, msg, "<expected_output>")
	})
}

func TestComputeSystemMessageWithSubordinates(t *testing.T) {
	t.Run("single subordinate without tools", func(t *testing.T) {
		subordinate := &mockAgent{
			name:        "helper-agent",
			description: "A helper agent",
			skills:      []tools.SkillCard{},
		}
		
		agent, err := New(
			createMockModel(),
			WithName("manager-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		assert.Contains(t, msg, "Your name is manager-agent")
		assert.Contains(t, msg, "<subordinates>")
		assert.Contains(t, msg, "Here are the members in your team:")
		assert.Contains(t, msg, " - Agent 1:")
		assert.Contains(t, msg, "   - Name: helper-agent")
		assert.Contains(t, msg, "   - Description: A helper agent")
		assert.Contains(t, msg, "</subordinates>")
		assert.Contains(t, msg, "<task_delegation>")
		assert.Contains(t, msg, "Depending on the nature of the user request")
		assert.Contains(t, msg, "</task_delegation>")
	})
	
	t.Run("subordinate with tools", func(t *testing.T) {
		subordinate := &mockAgent{
			name:        "tool-agent",
			description: "An agent with tools",
			skills: []tools.SkillCard{
				{Name: "calculate", Description: "Performs calculations"},
				{Name: "search", Description: "Searches for information"},
			},
		}
		
		agent, err := New(
			createMockModel(),
			WithName("supervisor-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		assert.Contains(t, msg, "   - Member tools:")
		assert.Contains(t, msg, "    - calculate:Performs calculations")
		assert.Contains(t, msg, "    - search:Searches for information")
	})
	
	t.Run("multiple subordinates", func(t *testing.T) {
		subordinate1 := &mockAgent{
			name:        "agent-1",
			description: "First agent",
			skills:      []tools.SkillCard{{Name: "skill1", Description: "First skill"}},
		}
		
		subordinate2 := &mockAgent{
			name:        "agent-2",
			description: "Second agent",
			skills:      []tools.SkillCard{{Name: "skill2", Description: "Second skill"}},
		}
		
		agent, err := New(
			createMockModel(),
			WithName("team-leader"),
			WithSubordinate(subordinate1),
			WithSubordinate(subordinate2),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		assert.Contains(t, msg, " - Agent 1:")
		assert.Contains(t, msg, "   - Name: agent-1")
		assert.Contains(t, msg, " - Agent 2:")
		assert.Contains(t, msg, "   - Name: agent-2")
		
		// Check that agents are separated by newlines
		lines := strings.Split(msg, "\n")
		var agentLines []int
		for i, line := range lines {
			if strings.Contains(line, " - Agent ") {
				agentLines = append(agentLines, i)
			}
		}
		assert.Len(t, agentLines, 2)
	})
	
	t.Run("subordinate with empty description", func(t *testing.T) {
		subordinate := &mockAgent{
			name:        "no-desc-agent",
			description: "",
			skills:      []tools.SkillCard{},
		}
		
		agent, err := New(
			createMockModel(),
			WithName("parent-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		assert.Contains(t, msg, "   - Name: no-desc-agent")
		// Description line should not be present when description is empty
		assert.NotContains(t, msg, "   - Description:")
	})
}

func TestComputeSystemMessageErrors(t *testing.T) {
	t.Run("subordinate name error", func(t *testing.T) {
		subordinate := &mockAgent{
			nameErr: assert.AnError,
		}
		
		agent, err := New(
			createMockModel(),
			WithName("error-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		_, err = agent.ComputeSystemMessage()
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
	
	t.Run("subordinate description error", func(t *testing.T) {
		subordinate := &mockAgent{
			name:    "test-agent",
			descErr: assert.AnError,
		}
		
		agent, err := New(
			createMockModel(),
			WithName("error-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		_, err = agent.ComputeSystemMessage()
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
	
	t.Run("subordinate skills error", func(t *testing.T) {
		subordinate := &mockAgent{
			name:        "test-agent",
			description: "test description",
			skillsErr:   assert.AnError,
		}
		
		agent, err := New(
			createMockModel(),
			WithName("error-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		_, err = agent.ComputeSystemMessage()
		assert.Error(t, err)
		assert.Equal(t, assert.AnError, err)
	})
}

func TestSystemMessageFormatting(t *testing.T) {
	t.Run("proper XML tag formatting", func(t *testing.T) {
		agent, err := New(
			createMockModel(),
			WithName("format-test"),
			WithGoal("Test goal"),
			WithInstructions("Test instructions"),
			WithExpectedOutput("Test output"),
		)
		require.NoError(t, err)
		
		msg, err := agent.ComputeSystemMessage()
		require.NoError(t, err)
		
		// Check proper XML formatting
		assert.Contains(t, msg, "<goal>\nTest goal\n</goal>")
		assert.Contains(t, msg, "<instructions>\nTest instructions\n</instructions>")
		assert.Contains(t, msg, "<expected_output>\nTest output\n</expected_output>")
		
		// Check sections are separated by newlines
		assert.Contains(t, msg, "Your name is format-test")
		assert.Contains(t, msg, "\n")
	})
}