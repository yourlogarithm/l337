package agent_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/tools"
)

func TestComputeSystemMessage_Empty(t *testing.T) {
	a, err := agent.New(nil)
	require.NoError(t, err)

	out, err := a.ComputeSystemMessage()
	require.NoError(t, err)

	assert.Empty(t, out)
}

func TestComputeSystemMessage_Name(t *testing.T) {
	a, err := agent.New(nil, agent.WithName("Alpha"))
	require.NoError(t, err)

	out, err := a.ComputeSystemMessage()
	require.NoError(t, err)

	expected := "Your name is Alpha."
	assert.Equal(t, expected, out)
}

func TestComputeSystemMessage_WithEverything(t *testing.T) {
	bravo, err := agent.New(nil,
		agent.WithName("Bravo"),
		agent.WithDescription("Bravo description"),
		agent.WithInstructions("Bravo instructions"),
		agent.WithExpectedOutput("Bravo expected output"),
		agent.WithTool(tools.New("bravo_tool_0", "Bravo tool description 0", nil)),
		agent.WithTool(tools.New("bravo_tool_1", "Bravo tool description 1", nil)),
	)
	require.NoError(t, err)

	charlie, err := agent.New(nil,
		agent.WithName("Charlie"),
		agent.WithDescription("Charlie description"),
		agent.WithInstructions("Charlie instructions"),
		agent.WithExpectedOutput("Charlie expected output"),
		agent.WithTool(tools.New("charlie_tool_0", "Charlie tool description 0", nil)),
	)
	require.NoError(t, err)

	alpha, err := agent.New(nil,
		agent.WithName("Alpha"),
		agent.WithDescription("Alpha description"),
		agent.WithInstructions("Alpha instructions"),
		agent.WithExpectedOutput("Alpha expected output"),
		agent.WithSubordinate(bravo),
		agent.WithSubordinate(charlie),
	)
	require.NoError(t, err)

	out, err := alpha.ComputeSystemMessage()
	require.NoError(t, err)

	expected, err := os.ReadFile("testdata/system_message_with_everything.golden")

	require.NoError(t, err)
	assert.Equal(t, string(expected), out)
}
