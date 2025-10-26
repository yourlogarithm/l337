package agent_test

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func TestAgent_RunStreaming_NilModelError(t *testing.T) {
	agentInstance, _ := agent.New(nil)
	_, _, err := agentInstance.RunStreamingWithParams(t.Context(), 1)
	assert.Equal(t, agent.ErrBuilderParams{Param: "model", Msg: "nil"}, err)
}

func TestAgent_RunStreaming_EmptyMessagesError(t *testing.T) {
	model := MockModel{}

	agentInstance, _ := agent.New(model.Wrap())

	_, _, err := agentInstance.RunStreamingWithParams(t.Context(), 1)
	assert.Equal(t, agent.ErrBuilderParams{Param: "messages", Msg: "at least one message is required"}, err)
}

func TestAgent_RunStreaming_ChatStreamingErrorSendAndReturn(t *testing.T) {
	model := MockModel{
		ChatStreamingFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error) {
			return nil, assert.AnError
		},
	}
	agentInstance, _ := agent.New(model.Wrap())

	_, channel, err := agentInstance.RunStreamingWithParams(t.Context(), 1, types.WithTextMessage(types.RoleUser, "Hello"))
	assert.NoError(t, err)

	chunk, err := channel.Next()
	assert.Equal(t, assert.AnError, err)
	assert.Empty(t, chunk)

	_, err = channel.Next()
	assert.Equal(t, io.EOF, err)
}

func TestAgent_RunStreaming_StreamErr(t *testing.T) {
	model := MockModel{
		ChatStreamingFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error) {
			ch := streaming.NewResponseChannel(1)
			go func() {
				defer ch.Close()
				ch.Send(&types.Response{Content: types.NewTextContent("Hello")})
				ch.SendErr(assert.AnError)
			}()
			return ch, nil
		},
	}
	agentInstance, _ := agent.New(model.Wrap())

	_, channel, err := agentInstance.RunStreamingWithParams(t.Context(), 1, types.WithTextMessage(types.RoleUser, "Hello"))
	assert.NoError(t, err)

	chunk, err := channel.Next()
	assert.NoError(t, err)
	assert.Equal(t, types.NewTextContent("Hello"), chunk.Content)

	_, err = channel.Next()
	assert.Equal(t, assert.AnError, err)

	_, err = channel.Next()
	assert.Equal(t, io.EOF, err)
}

func TestAgent_RunStreaming_AddChunkError(t *testing.T) {
	model := MockModel{
		ChatStreamingFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error) {
			ch := streaming.NewResponseChannel(1)
			go func() {
				defer ch.Close()
				ch.Send(&types.Response{ID: "a", Content: types.NewTextContent("Hello")})
				ch.Send(&types.Response{ID: "b", Content: types.NewTextContent("World")})
			}()
			return ch, nil
		},
	}
	agentInstance, _ := agent.New(model.Wrap())

	_, channel, err := agentInstance.RunStreamingWithParams(t.Context(), 1, types.WithTextMessage(types.RoleUser, "Hello"))
	assert.NoError(t, err)

	chunk, err := channel.Next()
	assert.NoError(t, err)
	assert.Equal(t, types.NewTextContent("Hello"), chunk.Content)

	_, err = channel.Next()
	assert.Error(t, err)

	errChunkAddition, ok := err.(streaming.ErrChunkAddition)
	assert.True(t, ok)

	assert.Equal(t, "a", errChunkAddition.Accumulator.ID)
	assert.Equal(t, types.NewTextContent("Hello"), errChunkAddition.Accumulator.Content)

	assert.Equal(t, "b", errChunkAddition.Chunk.ID)
	assert.Equal(t, types.NewTextContent("World"), errChunkAddition.Chunk.Content)

	_, err = channel.Next()
	assert.Equal(t, io.EOF, err)
}

func TestAgent_RunStreaming_AccumulateUntilEOF(t *testing.T) {
	model := MockModel{
		ChatStreamingFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error) {
			ch := streaming.NewResponseChannel(0)
			go func() {
				defer ch.Close()
				ch.Send(&types.Response{Content: types.NewTextContent("Hello")})
				ch.Send(&types.Response{Content: types.NewTextContent(" ")})
				ch.Send(&types.Response{Content: types.NewTextContent("World"), FinishReason: "stop"})
			}()
			return ch, nil
		},
	}
	agentInstance, _ := agent.New(model.Wrap())

	run, channel, err := agentInstance.RunStreamingWithParams(t.Context(), 0, types.WithTextMessage(types.RoleUser, "Hello"))
	assert.NoError(t, err)

	var accumulated string
	for {
		chunk, err := channel.Next()
		if err == io.EOF {
			break
		}
		assert.NoError(t, err)
		accumulated += chunk.Content.Text
	}

	assert.Equal(t, "Hello World", accumulated)
	assert.Equal(t, 2, len(run.Messages))
	assert.Equal(t, types.Message{Role: types.RoleUser, Content: types.NewTextContent("Hello")}, run.Messages[0])
	assert.Equal(t, types.Message{Role: types.RoleAssistant, Content: types.NewTextContent("Hello World")}, run.Messages[1])
	assert.Equal(t, 1, len(run.Metrics))
}

func TestAgent_RunStreaming_ContinueUntilNoToolsCalled(t *testing.T) {
	functionCalls := 0
	chatCalls := 0

	model := MockModel{
		ChatStreamingFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (streaming.ResponseChannel, error) {
			chatCalls++
			ch := streaming.NewResponseChannel(0)
			go func() {
				defer ch.Close()
				if chatCalls > 3 {
					ch.Send(&types.Response{Content: types.NewTextContent("Final response"), FinishReason: "stop"})
				} else {
					toolCalls := []types.ToolCall{
						{
							ID:        "test_tool#0",
							Name:      "test_tool",
							Arguments: "{\"input\":",
						},
						{
							Arguments: "\"some ",
						},
						{
							Arguments: " input\"",
						},
						{
							Arguments: "}",
						},
					}
					for _, toolCall := range toolCalls {
						ch.Send(&types.Response{ToolCalls: []types.ToolCall{toolCall}})
					}
					ch.Send(&types.Response{FinishReason: "tool_call"})
				}
			}()
			return ch, nil
		},
	}

	type arg struct {
		Input string `json:"input"`
	}

	callable := func(ctx context.Context, run *types.Run, arg arg) ([]types.Content, error) {
		functionCalls++
		return types.NewTextContent(arg.Input).AsSlice(), nil
	}

	tool, err := tools.NewWithArgs("test_tool", "A test tool", callable)
	require.NoError(t, err)

	agentInstance, _ := agent.New(model.Wrap(), agent.WithTool(tool))

	_, channel, err := agentInstance.RunStreamingWithParams(t.Context(), 0, types.WithTextMessage(types.RoleUser, "Hello"))
	assert.NoError(t, err)

	channel.Drain()

	assert.Equal(t, 3, functionCalls)
	assert.Equal(t, 4, chatCalls)
}
