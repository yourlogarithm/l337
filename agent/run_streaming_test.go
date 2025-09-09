package agent

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
)

func TestRunStreamingWithParams(t *testing.T) {
	t.Run("successful streaming with params", func(t *testing.T) {
		mockChannel := provider.NewResponseChannel(10)
		
		mockImpl := &mockModelImpl{
			streamResponse: mockChannel,
		}
		
		model := &provider.Model{
			Name:     "stream-model",
			Provider: "test",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("stream-agent"))
		require.NoError(t, err)
		
		runResponse, channel, err := agent.RunStreamingWithParams(
			context.Background(),
			10, // buffer size
			chat.WithMessage(chat.RoleUser, "Test streaming"),
		)
		
		require.NoError(t, err)
		assert.NotNil(t, runResponse)
		assert.NotNil(t, channel)
		assert.Len(t, runResponse.Messages, 2) // system + user messages
		assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
		assert.Contains(t, runResponse.Messages[0].Content, "Your name is stream-agent")
		assert.Equal(t, chat.RoleUser, runResponse.Messages[1].Role)
		assert.Equal(t, "Test streaming", runResponse.Messages[1].Content)
	})
	
	t.Run("error with empty messages", func(t *testing.T) {
		agent, err := New(createMockModel(), WithName("stream-agent"))
		require.NoError(t, err)
		
		runResponse, channel, err := agent.RunStreamingWithParams(context.Background(), 10)
		
		assert.Error(t, err)
		assert.Nil(t, runResponse)
		assert.Nil(t, channel)
		assert.Contains(t, err.Error(), "at least one message is required")
	})
}

func TestRunStreaming(t *testing.T) {
	t.Run("streaming without system message adds computed one", func(t *testing.T) {
		mockChannel := provider.NewResponseChannel(10)
		
		mockImpl := &mockModelImpl{
			streamResponse: mockChannel,
		}
		
		model := &provider.Model{
			Name:     "stream-model",
			Provider: "test",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("computed-stream-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "User message without system"},
			},
		}
		
		channel, err := agent.RunStreaming(context.Background(), runResponse, 10)
		
		require.NoError(t, err)
		assert.NotNil(t, channel)
		
		// System message should be computed and prepended
		assert.Len(t, runResponse.Messages, 2) // system + user
		assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
		assert.Contains(t, runResponse.Messages[0].Content, "Your name is computed-stream-agent")
	})
	
	t.Run("streaming with existing system message", func(t *testing.T) {
		mockChannel := provider.NewResponseChannel(5)
		
		mockImpl := &mockModelImpl{
			streamResponse: mockChannel,
		}
		
		model := &provider.Model{
			Name:     "stream-model",
			Provider: "test",
			Impl:     mockImpl,
		}
		
		agent, err := New(model, WithName("existing-stream-agent"))
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleSystem, Content: "Existing system message"},
				{Role: chat.RoleUser, Content: "User message"},
			},
		}
		
		channel, err := agent.RunStreaming(context.Background(), runResponse, 5)
		
		require.NoError(t, err)
		assert.NotNil(t, channel)
		
		// System message should remain unchanged
		assert.Len(t, runResponse.Messages, 2) // system + user
		assert.Equal(t, "Existing system message", runResponse.Messages[0].Content)
	})
	
	t.Run("streaming with system message computation error", func(t *testing.T) {
		// Create subordinate that will cause an error in system message computation
		subordinate := &mockAgent{
			name:    "error-subordinate",
			nameErr: assert.AnError,
		}
		
		agent, err := New(
			createMockModel(),
			WithName("error-stream-agent"),
			WithSubordinate(subordinate),
		)
		require.NoError(t, err)
		
		runResponse := &chat.RunResponse{
			SessionID: uuid.New(),
			Messages: []chat.Message{
				{Role: chat.RoleUser, Content: "This will cause system message error"},
			},
		}
		
		channel, err := agent.RunStreaming(context.Background(), runResponse, 10)
		
		assert.Error(t, err)
		assert.Nil(t, channel)
		assert.Equal(t, assert.AnError, err)
	})
}