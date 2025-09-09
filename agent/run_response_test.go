package agent

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/chat"
)

func TestBuildRunResponse(t *testing.T) {
	t.Run("successful build with messages", func(t *testing.T) {
		sessionID := uuid.New()

		runResponse, err := BuildRunResponse(
			chat.WithSessionID(sessionID),
			chat.WithMessage(chat.RoleUser, "Hello world"),
		)

		require.NoError(t, err)
		assert.NotNil(t, runResponse)
		assert.Equal(t, sessionID, runResponse.SessionID)
		assert.Len(t, runResponse.Messages, 1)
		assert.Equal(t, chat.RoleUser, runResponse.Messages[0].Role)
		assert.Equal(t, "Hello world", runResponse.Messages[0].Content)
		assert.NotNil(t, runResponse.Metrics)
	})

	t.Run("successful build with multiple messages", func(t *testing.T) {
		runResponse, err := BuildRunResponse(
			chat.WithMessage(chat.RoleSystem, "You are a helpful assistant"),
			chat.WithMessage(chat.RoleUser, "Hello"),
			chat.WithMessage(chat.RoleAssistant, "Hi there!"),
		)

		require.NoError(t, err)
		assert.Len(t, runResponse.Messages, 3)

		assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
		assert.Equal(t, "You are a helpful assistant", runResponse.Messages[0].Content)

		assert.Equal(t, chat.RoleUser, runResponse.Messages[1].Role)
		assert.Equal(t, "Hello", runResponse.Messages[1].Content)

		assert.Equal(t, chat.RoleAssistant, runResponse.Messages[2].Role)
		assert.Equal(t, "Hi there!", runResponse.Messages[2].Content)
	})

	t.Run("empty messages should return error", func(t *testing.T) {
		runResponse, err := BuildRunResponse()

		assert.Error(t, err)
		assert.Nil(t, runResponse)
		assert.Contains(t, err.Error(), "at least one message is required")
	})

	t.Run("only session ID without messages should return error", func(t *testing.T) {
		sessionID := uuid.New()

		runResponse, err := BuildRunResponse(
			chat.WithSessionID(sessionID),
		)

		assert.Error(t, err)
		assert.Nil(t, runResponse)
		assert.Contains(t, err.Error(), "at least one message is required")
	})

	t.Run("builds with default session ID when not provided", func(t *testing.T) {
		runResponse, err := BuildRunResponse(
			chat.WithMessage(chat.RoleUser, "Test message"),
		)

		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, runResponse.SessionID)
		assert.NotNil(t, runResponse)
		assert.Len(t, runResponse.Messages, 1)
	})

	t.Run("handles multiple parameters of same type", func(t *testing.T) {
		runResponse, err := BuildRunResponse(
			chat.WithMessage(chat.RoleUser, "First message"),
			chat.WithMessage(chat.RoleUser, "Second message"),
			chat.WithMessage(chat.RoleAssistant, "Response message"),
		)

		require.NoError(t, err)
		assert.Len(t, runResponse.Messages, 3)
		assert.Equal(t, "First message", runResponse.Messages[0].Content)
		assert.Equal(t, "Second message", runResponse.Messages[1].Content)
		assert.Equal(t, "Response message", runResponse.Messages[2].Content)
	})
}
