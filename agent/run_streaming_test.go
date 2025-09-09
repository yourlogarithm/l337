package agent

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/provider"
)

func TestRunStreamingStopReason(t *testing.T) {
	mockChannel := provider.NewResponseChannel(10)
	now := time.Now()

	go func() {
		defer mockChannel.Close()
		mockChannel.Send(&provider.Response{ID: "resp-1", Created: now})

		mockChannel.Send(&provider.Response{Reasoning: "Reasonable "})
		mockChannel.Send(&provider.Response{Reasoning: "reasoning."})

		mockChannel.Send(&provider.Response{Content: "Contentful"})
		mockChannel.Send(&provider.Response{Content: " content."})

		mockChannel.Send(&provider.Response{FinishReason: "stop"})

		mockChannel.Send(&provider.Response{Metrics: metrics.Metrics{Timestamp: time.Now(), TotalTokens: 10, TotalDuration: 500 * time.Millisecond}})
	}()

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

	channel.WaitUntilFinished()

	assert.NotNil(t, runResponse)
	assert.NotNil(t, channel)
	assert.Len(t, runResponse.Messages, 3) // system + user messages
	assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
	assert.Contains(t, runResponse.Messages[0].Content, "Your name is stream-agent")
	assert.Equal(t, chat.RoleUser, runResponse.Messages[1].Role)
	assert.Equal(t, "Test streaming", runResponse.Messages[1].Content)
	assert.Equal(t, chat.RoleAssistant, runResponse.Messages[2].Role)
	assert.Equal(t, "Reasonable reasoning.", runResponse.Messages[2].Reasoning)
	assert.Equal(t, "Contentful content.", runResponse.Messages[2].Content)

	assert.NotEmpty(t, runResponse.Metrics)
	metrics, ok := runResponse.Metrics[agent.id]
	assert.True(t, ok)
	assert.NotNil(t, metrics)

	assert.Equal(t, 1, len(metrics))
	assert.Equal(t, uint(10), metrics[0].TotalTokens)
	assert.Equal(t, 500*time.Millisecond, metrics[0].TotalDuration)
}

func TestRunStreamingToolCall(t *testing.T) {
	mockChannel := provider.NewResponseChannel(10)
	now := time.Now()

	go func() {
		defer mockChannel.Close()
		mockChannel.Send(&provider.Response{ID: "resp-1", Created: now})

		mockChannel.Send(&provider.Response{})
	}()

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

	channel.WaitUntilFinished()

	assert.NotNil(t, runResponse)
	assert.NotNil(t, channel)
	assert.Len(t, runResponse.Messages, 3) // system + user messages
	assert.Equal(t, chat.RoleSystem, runResponse.Messages[0].Role)
	assert.Contains(t, runResponse.Messages[0].Content, "Your name is stream-agent")
	assert.Equal(t, chat.RoleUser, runResponse.Messages[1].Role)
	assert.Equal(t, "Test streaming", runResponse.Messages[1].Content)
	assert.Equal(t, chat.RoleAssistant, runResponse.Messages[2].Role)
	assert.Equal(t, "Reasonable reasoning.", runResponse.Messages[2].Reasoning)
	assert.Equal(t, "Contentful content.", runResponse.Messages[2].Content)

	assert.NotEmpty(t, runResponse.Metrics)
	metrics, ok := runResponse.Metrics[agent.id]
	assert.True(t, ok)
	assert.NotNil(t, metrics)

	assert.Equal(t, 1, len(metrics))
	assert.Equal(t, uint(10), metrics[0].TotalTokens)
	assert.Equal(t, 500*time.Millisecond, metrics[0].TotalDuration)
}
