package agent

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/provider"
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
		
		var delegateToolFound bool
		for _, skill := range skills {
			if skill.Name == "delegate_task" {
				delegateToolFound = true
				assert.Equal(t, "Delegates the task to one or more subordinates", skill.Description)
				break
			}
		}
		assert.True(t, delegateToolFound, "delegate_task tool should be automatically added when subordinates exist")
	})
	
	t.Run("no delegate tool without subordinates", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("solo-agent"))
		require.NoError(t, err)
		
		skills, err := agent.Skills()
		require.NoError(t, err)
		
		// Should not have delegate_task tool without subordinates
		for _, skill := range skills {
			assert.NotEqual(t, "delegate_task", skill.Name)
		}
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

func TestDelegateTaskTool(t *testing.T) {
	t.Run("successful delegation to single subordinate", func(t *testing.T) {
		// Create a mock subordinate that returns a specific response
		subordinate := &mockAgentWithRunFunc{
			mockAgent: mockAgent{
				name:        "worker",
				description: "A worker agent",
				skills:      []tools.SkillCard{},
			},
			runFunc: func(ctx context.Context, runResponse *chat.RunResponse) error {
				// Simulate the subordinate responding to the delegated task
				runResponse.Messages = append(runResponse.Messages, chat.Message{
					Role:    chat.RoleAssistant,
					Content: "Task completed successfully by worker",
				})
				return nil
			},
		}
		
		// Create manager with the subordinate
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content: "I'll delegate this task",
				ToolCalls: []chat.ToolCall{
					{
						ID:        "call_1", 
						Name:      "delegate_task", 
						Arguments: `{"names": ["worker"], "expected_output": "Complete the assigned work"}`,
					},
				},
				FinishReason: "tool_calls",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			},
		}
		
		var callCount int
		mockImpl.chatFunc = func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			callCount++
			if callCount == 1 {
				return mockImpl.chatResponse, nil
			}
			return provider.Response{
				Content:      "The task has been completed by the subordinate",
				FinishReason: "stop",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			}, nil
		}
		
		model := &provider.Model{
			Name:     "manager-model",
			Provider: "test",
			Impl:     mockImpl,
		}
		
		agent, err := New(
			model,
			WithName("manager"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "Please delegate this task to a worker"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		require.NoError(t, err)
		
		// Check that delegation occurred and subordinate response was included
		toolResponseFound := false
		for _, msg := range runResponse.Messages {
			if msg.Role == chat.RoleTool && msg.Name == "delegate_task" {
				toolResponseFound = true
				assert.Contains(t, msg.Content, "(worker) Response: Task completed successfully by worker")
			}
		}
		assert.True(t, toolResponseFound, "Tool response with subordinate output should be present")
	})
	
	t.Run("delegation to multiple subordinates", func(t *testing.T) {
		// Create multiple subordinates
		subordinate1 := &mockAgentWithRunFunc{
			mockAgent: mockAgent{
				name: "worker1",
				skills: []tools.SkillCard{},
			},
			runFunc: func(ctx context.Context, runResponse *chat.RunResponse) error {
				runResponse.Messages = append(runResponse.Messages, chat.Message{
					Role:    chat.RoleAssistant,
					Content: "Worker1 completed the task",
				})
				return nil
			},
		}
		
		subordinate2 := &mockAgentWithRunFunc{
			mockAgent: mockAgent{
				name: "worker2", 
				skills: []tools.SkillCard{},
			},
			runFunc: func(ctx context.Context, runResponse *chat.RunResponse) error {
				runResponse.Messages = append(runResponse.Messages, chat.Message{
					Role:    chat.RoleAssistant,
					Content: "Worker2 completed the task",
				})
				return nil
			},
		}
		
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content: "I'll delegate to multiple workers",
				ToolCalls: []chat.ToolCall{
					{
						ID:        "call_1", 
						Name:      "delegate_task", 
						Arguments: `{"names": ["worker1", "worker2"], "expected_output": "Complete your part of the work"}`,
					},
				},
				FinishReason: "tool_calls",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			},
		}
		
		var callCount int
		mockImpl.chatFunc = func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			callCount++
			if callCount == 1 {
				return mockImpl.chatResponse, nil
			}
			return provider.Response{
				Content:      "Both workers have completed their tasks",
				FinishReason: "stop",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			}, nil
		}
		
		model := &provider.Model{
			Name:     "multi-manager-model",
			Provider: "test",
			Impl:     mockImpl,
		}
		
		agent, err := New(
			model,
			WithName("multi-manager"),
			WithSubordinate(subordinate1),
			WithSubordinate(subordinate2),
		)
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "Delegate work to both workers"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		require.NoError(t, err)
		
		// Find the tool response and check both subordinates responded
		var toolResponse string
		for _, msg := range runResponse.Messages {
			if msg.Role == chat.RoleTool && msg.Name == "delegate_task" {
				toolResponse = msg.Content
				break
			}
		}
		
		assert.Contains(t, toolResponse, "(worker1) Response: Worker1 completed the task")
		assert.Contains(t, toolResponse, "(worker2) Response: Worker2 completed the task")
	})
	
	t.Run("delegation error handling", func(t *testing.T) {
		// Create subordinate that returns an error
		subordinate := &mockAgentWithRunFunc{
			mockAgent: mockAgent{
				name: "error-worker",
				skills: []tools.SkillCard{},
			},
			runFunc: func(ctx context.Context, runResponse *chat.RunResponse) error {
				return errors.New("subordinate task failed")
			},
		}
		
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content: "I'll delegate this task",
				ToolCalls: []chat.ToolCall{
					{
						ID:        "call_1", 
						Name:      "delegate_task", 
						Arguments: `{"names": ["error-worker"], "expected_output": "Complete the work"}`,
					},
				},
				FinishReason: "tool_calls",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			},
		}
		
		var callCount int
		mockImpl.chatFunc = func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			callCount++
			if callCount == 1 {
				return mockImpl.chatResponse, nil
			}
			return provider.Response{
				Content:      "I handled the subordinate error",
				FinishReason: "stop",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			}, nil
		}
		
		model := &provider.Model{
			Name:     "error-manager-model",
			Provider: "test",
			Impl:     mockImpl,
		}
		
		agent, err := New(
			model,
			WithName("error-manager"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "Delegate to error worker"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		require.NoError(t, err)
		
		// Find the tool response and check error was reported
		var toolResponse string
		for _, msg := range runResponse.Messages {
			if msg.Role == chat.RoleTool && msg.Name == "delegate_task" {
				toolResponse = msg.Content
				break
			}
		}
		
		assert.Contains(t, toolResponse, "(error-worker) Error: subordinate task failed")
	})
	
	t.Run("delegation with invalid parameters", func(t *testing.T) {
		subordinate := &mockAgent{
			name: "worker",
			skills: []tools.SkillCard{},
		}
		
		agent, err := New(
			createMockModel(),
			WithName("invalid-manager"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		// Get the delegate_task tool
		var delegateTool *tools.Tool
		for _, tool := range agent.tools {
			if tool.SkillCard.Name == "delegate_task" {
				delegateTool = &tool
				break
			}
		}
		require.NotNil(t, delegateTool)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages:  []chat.Message{},
			Metrics:   make(map[uuid.UUID][]metrics.Metrics),
		}
		
		// Test with empty names array
		result, err := delegateTool.Callable(
			context.Background(),
			runResponse,
			`{"names": [], "expected_output": "Complete task"}`,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no subordinates specified")
		assert.Empty(t, result)
		
		// Test with empty expected output
		result, err = delegateTool.Callable(
			context.Background(),
			runResponse,
			`{"names": ["worker"], "expected_output": ""}`,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no expected output specified")
		assert.Empty(t, result)
	})
}

// mockAgentWithRunFunc extends mockAgent to allow custom run function
type mockAgentWithRunFunc struct {
	mockAgent
	runFunc func(context.Context, *chat.RunResponse) error
}

func (m *mockAgentWithRunFunc) Run(ctx context.Context, runResponse *chat.RunResponse) error {
	if m.runFunc != nil {
		return m.runFunc(ctx, runResponse)
	}
	return m.mockAgent.Run(ctx, runResponse)
}