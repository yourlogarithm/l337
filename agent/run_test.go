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

func TestRunWithParams(t *testing.T) {
	t.Run("successful run with minimal params", func(t *testing.T) {
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content:      "Test response from model",
				FinishReason: "stop",
				Metrics: metrics.Metrics{
					Timestamp:        time.Now(),
					TotalTokens:      15,
					PromptTokens:     10,
					CompletionTokens: 5,
					TotalDuration:    100 * time.Millisecond,
				},
			},
		}
		
		model := &provider.Model{
			Name:     "test-model",
			Provider: "test-provider",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("test-agent"))
		require.NoError(t, err)
		
		runResponse, err := agent.RunWithParams(
			context.Background(),
			chat.WithMessage(chat.RoleUser, "Hello, test agent!"),
		)
		
		require.NoError(t, err)
		assert.NotNil(t, runResponse)
		
		// Check that messages include system, user, and assistant messages
		assert.Len(t, runResponse.Messages, 3)
		assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
		assert.Equal(t, chat.RoleUser, runResponse.Messages[1].Role)
		assert.Equal(t, chat.RoleAssistant, runResponse.Messages[2].Role)
		assert.Equal(t, "Test response from model", runResponse.Messages[2].Content)
	})
	
	t.Run("empty messages should return error", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("test-agent"))
		require.NoError(t, err)
		
		runResponse, err := agent.RunWithParams(context.Background())
		
		assert.Error(t, err)
		assert.Nil(t, runResponse)
		assert.Contains(t, err.Error(), "at least one message is required")
	})
}

func TestRun(t *testing.T) {
	t.Run("run with existing system message", func(t *testing.T) {
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content:      "Response with existing system message",
				FinishReason: "stop",
				Metrics: metrics.Metrics{
					Timestamp: time.Now(),
				},
			},
		}
		
		model := &provider.Model{
			Name:     "test-model",
			Provider: "test-provider",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("test-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleSystem, Content: "Existing system message"},
				{Role: chat.RoleUser, Content: "User message"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		
		require.NoError(t, err)
		
		// System message should remain the existing one, not be replaced
		assert.Equal(t, "Existing system message", runResponse.Messages[0].Content)
		assert.Len(t, runResponse.Messages, 3) // system, user, assistant
		assert.Equal(t, "Response with existing system message", runResponse.Messages[2].Content)
	})
	
	t.Run("run without system message adds computed one", func(t *testing.T) {
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content:      "Response with computed system message",
				FinishReason: "stop",
				Metrics: metrics.Metrics{
					Timestamp: time.Now(),
				},
			},
		}
		
		model := &provider.Model{
			Name:     "test-model",
			Provider: "test-provider",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("computed-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "User message without system"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		
		require.NoError(t, err)
		
		// System message should be computed and prepended
		assert.Len(t, runResponse.Messages, 3) // system, user, assistant
		assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
		assert.Contains(t, runResponse.Messages[0].Content, "Your name is computed-agent")
	})
	
	t.Run("model error is propagated", func(t *testing.T) {
		expectedErr := errors.New("model failed")
		mockImpl := &mockModelImpl{
			chatError: expectedErr,
		}
		
		model := &provider.Model{
			Name:     "error-model",
			Provider: "test-provider",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("error-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "This will fail"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestRunWithTools(t *testing.T) {
	t.Run("tool call execution", func(t *testing.T) {
		// First call returns tool call, second call returns final response
		var callCount int
		mockImpl := &mockModelImpl{}
		mockImpl.chatResponse = provider.Response{
			Content: "I need to use a tool",
			ToolCalls: []chat.ToolCall{
				{ID: "call_1", Name: "test_tool", Arguments: `{"value": "test"}`},
			},
			FinishReason: "tool_calls",
			Metrics: metrics.Metrics{Timestamp: time.Now()},
		}
		
		// Mock the Chat method to return different responses on successive calls
		mockImpl.chatFunc = func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			callCount++
			if callCount == 1 {
				return provider.Response{
					Content: "I need to use a tool",
					ToolCalls: []chat.ToolCall{
						{ID: "call_1", Name: "test_tool", Arguments: `{"value": "test"}`},
					},
					FinishReason: "tool_calls",
					Metrics: metrics.Metrics{Timestamp: time.Now()},
				}, nil
			}
			return provider.Response{
				Content:      "Tool result processed: test_result",
				FinishReason: "stop",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			}, nil
		}
		
		model := &provider.Model{
			Name:     "tool-model",
			Provider: "test-provider",
			Impl:     mockImpl,
		}
		
		// Create a test tool
		type TestArgs struct {
			Value string `json:"value"`
		}
		
		testTool, err := tools.NewToolWithArgs("test_tool", "A test tool",
			func(ctx context.Context, response *chat.RunResponse, args TestArgs) (string, error) {
				return "test_result", nil
			})
		require.NoError(t, err)
		
		agent, err := New(model, WithName("tool-agent"), WithTool(testTool))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "Please use the test tool"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		err = agent.Run(context.Background(), runResponse)
		
		require.NoError(t, err)
		
		// Should have: system, user, assistant (with tool call), tool result, final assistant
		assert.Len(t, runResponse.Messages, 5)
		assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
		assert.Equal(t, chat.RoleUser, runResponse.Messages[1].Role)
		assert.Equal(t, chat.RoleAssistant, runResponse.Messages[2].Role)
		assert.Len(t, runResponse.Messages[2].ToolCalls, 1)
		assert.Equal(t, chat.RoleTool, runResponse.Messages[3].Role)
		assert.Equal(t, "test_result", runResponse.Messages[3].Content)
		assert.Equal(t, chat.RoleAssistant, runResponse.Messages[4].Role)
		assert.Equal(t, "Tool result processed: test_result", runResponse.Messages[4].Content)
	})
	
	t.Run("tool not found error", func(t *testing.T) {
		mockImpl := &mockModelImpl{
			chatResponse: provider.Response{
				Content: "I need to use a nonexistent tool",
				ToolCalls: []chat.ToolCall{
					{ID: "call_1", Name: "nonexistent_tool", Arguments: `{"value": "test"}`},
				},
				FinishReason: "tool_calls",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			},
		}
		
		model := &provider.Model{
			Name:     "tool-model",
			Provider: "test-provider",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("tool-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "Use nonexistent tool"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		// Mock the Chat method to return final response on second call
		callCount := 0
		mockImpl.chatFunc = func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			callCount++
			if callCount == 1 {
				return mockImpl.chatResponse, nil
			}
			return provider.Response{
				Content:      "I handled the tool error",
				FinishReason: "stop",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			}, nil
		}
		
		err = agent.Run(context.Background(), runResponse)
		
		require.NoError(t, err)
		
		// Check that tool error message was added
		assert.Contains(t, runResponse.Messages[3].Content, "error: tool not found: nonexistent_tool")
		assert.Equal(t, chat.RoleTool, runResponse.Messages[3].Role)
	})
	
	t.Run("tool execution error", func(t *testing.T) {
		mockImpl := &mockModelImpl{}
		
		model := &provider.Model{
			Name:     "tool-model",
			Provider: "test-provider", 
			Impl:     mockImpl,
		}
		
		// Create a tool that returns an error
		type TestArgs struct {
			Value string `json:"value"`
		}
		
		errorTool, err := tools.NewToolWithArgs("error_tool", "A tool that errors",
			func(ctx context.Context, response *chat.RunResponse, args TestArgs) (string, error) {
				return "", errors.New("tool execution failed")
			})
		require.NoError(t, err)
		
		agent, err := New(model, WithName("error-tool-agent"), WithTool(errorTool))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "Use error tool"},
			},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		// Mock the Chat method to return tool call then final response
		callCount := 0
		mockImpl.chatFunc = func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			callCount++
			if callCount == 1 {
				return provider.Response{
					Content: "I'll use the error tool",
					ToolCalls: []chat.ToolCall{
						{ID: "call_1", Name: "error_tool", Arguments: `{"value": "test"}`},
					},
					FinishReason: "tool_calls",
					Metrics: metrics.Metrics{Timestamp: time.Now()},
				}, nil
			}
			return provider.Response{
				Content:      "I handled the tool error",
				FinishReason: "stop",
				Metrics: metrics.Metrics{Timestamp: time.Now()},
			}, nil
		}
		
		err = agent.Run(context.Background(), runResponse)
		
		require.NoError(t, err)
		
		// Check that tool error message was added
		toolMessage := runResponse.Messages[3]
		assert.Equal(t, chat.RoleTool, toolMessage.Role)
		assert.Contains(t, toolMessage.Content, "error: tool execution failed")
		assert.True(t, toolMessage.IsErr)
	})
}

func TestHandleResponse(t *testing.T) {
	t.Run("nil response error", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("test-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			Messages: []chat.Message{},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		toolsCalled, err := agent.handleResponse(context.Background(), runResponse, nil)
		
		assert.Error(t, err)
		assert.False(t, toolsCalled)
		assert.Contains(t, err.Error(), "model response is nil")
	})
	
	t.Run("empty finish reason error", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("test-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			Messages: []chat.Message{},
			Metrics: make(map[uuid.UUID][]metrics.Metrics),
		}
		
		response := &provider.Response{
			Content: "test content",
			// FinishReason is empty
		}
		
		toolsCalled, err := agent.handleResponse(context.Background(), runResponse, response)
		
		assert.Error(t, err)
		assert.False(t, toolsCalled)
		assert.Contains(t, err.Error(), "response has no finish reason")
	})
	
	t.Run("successful response without tools", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("test-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages:  []chat.Message{},
			Metrics:   make(map[uuid.UUID][]metrics.Metrics),
		}
		
		response := &provider.Response{
			Content:      "test response",
			FinishReason: "stop",
			Metrics: metrics.Metrics{
				Timestamp: time.Now(),
			},
		}
		
		toolsCalled, err := agent.handleResponse(context.Background(), runResponse, response)
		
		require.NoError(t, err)
		assert.False(t, toolsCalled)
		
		// Check assistant message was added
		assert.Len(t, runResponse.Messages, 1)
		assert.Equal(t, chat.RoleAssistant, runResponse.Messages[0].Role)
		assert.Equal(t, "test response", runResponse.Messages[0].Content)
		
		// Check metrics were recorded
		assert.Len(t, runResponse.Metrics[agent.id], 1)
		assert.Equal(t, runResponse.SessionID, runResponse.Metrics[agent.id][0].SessionID)
	})
}