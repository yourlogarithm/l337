package agent_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/tools"
)

func TestNewAgent_WithMinimalConfiguration(t *testing.T) {
	agentInstance, err := agent.New(nil, agent.WithName("TestAgent"))
	require.NoError(t, err)

	name, err := agentInstance.Name()
	require.NoError(t, err)

	assert.Equal(t, "TestAgent", name)
	description, err := agentInstance.Description()
	require.NoError(t, err)
	assert.Empty(t, description)

	tools, err := agentInstance.Tools()
	require.NoError(t, err)
	assert.Empty(t, tools)
}

func TestNewAgent_WithAllOptions(t *testing.T) {
	tool := tools.New("test_tool", "A test tool", nil)
	subordinate, err := agent.New(nil, agent.WithName("Subordinate"))
	require.NoError(t, err)

	agentInstance, err := agent.New(nil,
		agent.WithID(uuid.New()),
		agent.WithName("TestAgent"),
		agent.WithDescription("A test agent"),
		agent.WithInstructions("Follow these instructions"),
		agent.WithExpectedOutput("Expected output"),
		agent.WithTool(tool),
		agent.WithSubordinate(subordinate),
	)
	require.NoError(t, err)

	name, err := agentInstance.Name()
	require.NoError(t, err)
	assert.Equal(t, "TestAgent", name)

	description, err := agentInstance.Description()
	require.NoError(t, err)
	assert.Equal(t, "A test agent", description)

	agentTools, err := agentInstance.Tools()
	require.NoError(t, err)

	assert.Len(t, agentTools, 2) // Including the delegate_task tool

	toolMap := make(map[string]tools.Tool)
	for _, t := range agentTools {
		toolMap[t.Name] = t
	}

	assert.Equal(t, "test_tool", agentTools[0].Name)
	assert.Equal(t, "A test tool", agentTools[0].Description)

	assert.Equal(t, agent.DELEGATE_TASK_TOOL_NAME, agentTools[1].Name)
	assert.Equal(t, agent.DELEGATE_TASK_TOOL_DESC, agentTools[1].Description)
}
