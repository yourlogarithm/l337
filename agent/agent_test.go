package agent

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

// mockModelImpl implements provider.ModelImpl for testing
type mockModelImpl struct {
	chatResponse   provider.Response
	chatError      error
	streamResponse provider.ResponseChannel
	streamError    error
	chatFunc       func(context.Context, *provider.Request, *chat.Options) (provider.Response, error)
}

func (m *mockModelImpl) Chat(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
	if m.chatFunc != nil {
		return m.chatFunc(ctx, req, opts)
	}
	return m.chatResponse, m.chatError
}

func (m *mockModelImpl) ChatStreaming(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.ResponseChannel, error) {
	return m.streamResponse, m.streamError
}

func createMockModel() *provider.Model {
	return &provider.Model{
		Name:     "test-model",
		Provider: "test-provider",
		Impl: &mockModelImpl{
			chatResponse: provider.Response{
				Content:      "Test response",
				FinishReason: "stop",
			},
		},
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		model         *provider.Model
		options       []AgentOption
		expectedError string
	}{
		{
			name:  "successful creation with minimal params",
			model: createMockModel(),
		},
		{
			name:          "nil model should return error",
			model:         nil,
			expectedError: "error building agent: model: model is required",
		},
		{
			name:  "successful creation with options",
			model: createMockModel(),
			options: []AgentOption{
				WithName("test-agent"),
				WithDescription("A test agent"),
				WithRole("assistant"),
				WithInstructions("Follow test instructions"),
				WithGoal("Achieve test goals"),
				WithExpectedOutput("Provide test output"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agent, err := New(tt.model, tt.options...)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, agent)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, agent)

				// Verify agent has basic properties
				assert.NotEqual(t, uuid.Nil, agent.id)
				assert.Equal(t, agent.name, agent.id.String())
				assert.Equal(t, tt.model, agent.model)
				assert.NotNil(t, agent.tools)
				assert.Empty(t, agent.tools)
				assert.NotNil(t, agent.subordinates)
				assert.Empty(t, agent.subordinates)
			}
		})
	}
}

func TestAgentBasicMethods(t *testing.T) {
	agent, err := New(createMockModel(), WithName("test-agent"), WithDescription("test description"))
	require.NoError(t, err)

	t.Run("Name", func(t *testing.T) {
		name, err := agent.Name()
		assert.NoError(t, err)
		assert.Equal(t, "test-agent", name)
	})

	t.Run("Description", func(t *testing.T) {
		description, err := agent.Description()
		assert.NoError(t, err)
		assert.Equal(t, "test description", description)
	})

	t.Run("Skills", func(t *testing.T) {
		skills, err := agent.Skills()
		assert.NoError(t, err)
		assert.NotNil(t, skills)
		assert.Empty(t, skills) // No tools added yet
	})
}

func TestWithOptions(t *testing.T) {
	testID := uuid.New()

	agent, err := New(
		createMockModel(),
		WithID(testID),
		WithName("custom-agent"),
		WithRole("test-role"),
		WithDescription("custom description"),
		WithInstructions("custom instructions"),
		WithGoal("custom goal"),
		WithExpectedOutput("custom output"),
	)
	require.NoError(t, err)

	assert.Equal(t, testID, agent.id)
	assert.Equal(t, "custom-agent", agent.name)
	assert.Equal(t, "test-role", agent.role)
	assert.Equal(t, "custom description", agent.description)
	assert.Equal(t, "custom instructions", agent.instructions)
	assert.Equal(t, "custom goal", agent.goal)
	assert.Equal(t, "custom output", agent.expectedOutput)
}

func TestWithID(t *testing.T) {
	t.Run("valid UUID", func(t *testing.T) {
		testID := uuid.New()
		agent, err := New(createMockModel(), WithID(testID))
		require.NoError(t, err)
		assert.Equal(t, testID, agent.id)
	})

	t.Run("nil UUID generates new UUID", func(t *testing.T) {
		agent, err := New(createMockModel(), WithID(uuid.Nil))
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, agent.id)
	})
}

func TestWithName(t *testing.T) {
	t.Run("valid name", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("test-agent"))
		require.NoError(t, err)
		assert.Equal(t, "test-agent", agent.name)
	})

	t.Run("empty name keeps default", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName(""))
		require.NoError(t, err)
		assert.NotEmpty(t, agent.name)
		assert.Equal(t, agent.name, agent.id.String()) // Should keep the default UUID-based name
	})
}

func TestWithTool(t *testing.T) {
	callable := func(ctx context.Context) (string, error) {
		return "test result", nil
	}

	testTool := tools.NewTool("test-tool", "A test tool", callable, tools.WithTags("test", "tag"), tools.WithExamples("Example", "Usage"))

	agent, err := New(createMockModel(), WithTool(testTool))
	require.NoError(t, err)

	skills, err := agent.Skills()
	require.NoError(t, err)
	assert.Len(t, skills, 1)
	assert.Equal(t, "test-tool", skills[0].Name)
	assert.Equal(t, "A test tool", skills[0].Description)
	assert.Equal(t, []string{"test", "tag"}, skills[0].Tags)
	assert.Equal(t, []string{"Example", "Usage"}, skills[0].Examples)
}

func TestWithChatOptions(t *testing.T) {
	temp := 0.7
	chatOptions := chat.Options{
		Temperature: &temp,
		MaxTokens:   1000,
	}

	agent, err := New(createMockModel(), WithChatOptions(chatOptions))
	require.NoError(t, err)

	assert.Equal(t, chatOptions, agent.chatOptions)
}
