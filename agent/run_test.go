package agent_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func TestAgent_Run_NilModel(t *testing.T) {
	agentInstance, err := agent.New(nil, agent.WithName("TestAgent"))
	require.NoError(t, err)

	_, err = agentInstance.RunWithParams(t.Context(), types.WithTextMessage(types.RoleUser, "Hello"))
	assert.Error(t, err)
	assert.Equal(t, agent.ErrBuilderParams{Param: "model", Msg: "nil"}, err)
}

func TestAgent_Run_EmptyMessages(t *testing.T) {
	model := MockModel{}

	agentInstance, err := agent.New(model.Wrap(), agent.WithName("TestAgent"))
	require.NoError(t, err)

	_, err = agentInstance.RunWithParams(t.Context())
	assert.Error(t, err)
	assert.Equal(t, agent.ErrBuilderParams{Param: "messages", Msg: "at least one message is required"}, err)
}

func TestAgent_Run_SystemMessageNotAppendedIfEmpty(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			assert.Equal(t, 1, len(messages))
			assert.Equal(t, types.Message{Role: types.RoleUser, Content: types.NewTextContent("Hello")}, messages[0])
			return types.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(model.Wrap())

	_, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)

	assert.NoError(t, err)
}

func TestAgent_Run_SystemMessageAppendedIfNotProvided(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			assert.Equal(t, 2, len(messages))
			assert.Equal(t, types.Message{Role: types.RoleSystem, Content: types.NewTextContent("Your name is TestAgent.\nA test agent")}, messages[0])
			assert.Equal(t, types.Message{Role: types.RoleUser, Content: types.NewTextContent("Hello")}, messages[1])
			return types.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithName("TestAgent"),
		agent.WithDescription("A test agent"),
	)

	_, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)
	assert.NoError(t, err)
}

func TestAgent_Run_SystemMessageNotAppendedIfProvided(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			assert.Equal(t, 2, len(messages))
			assert.Equal(t, types.Message{Role: types.RoleSystem, Content: types.NewTextContent("Custom system message")}, messages[0])
			assert.Equal(t, types.Message{Role: types.RoleUser, Content: types.NewTextContent("Hello")}, messages[1])
			return types.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithName("TestAgent"),
		agent.WithDescription("A test agent"),
	)

	_, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleSystem, "Custom system message"),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)
	assert.NoError(t, err)
}

func TestAgent_Run_ToolsPassed(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, toolsSlice []tools.Tool, opts *types.Options) (types.Response, error) {
			assert.Equal(t, 2, len(toolsSlice)) // Including the delegate_task tool

			toolMap := make(map[string]tools.Tool)
			for _, t := range toolsSlice {
				toolMap[t.Name] = t
			}

			assert.Equal(t, "test_tool", toolMap["test_tool"].Name)
			assert.Equal(t, "A test tool", toolMap["test_tool"].Description)

			assert.Equal(t, agent.DELEGATE_TASK_TOOL_NAME, toolMap[agent.DELEGATE_TASK_TOOL_NAME].Name)
			assert.Equal(t, agent.DELEGATE_TASK_TOOL_DESC, toolMap[agent.DELEGATE_TASK_TOOL_NAME].Description)

			return types.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	tool := tools.New("test_tool", "A test tool", nil)

	subordinate, _ := agent.New(nil, agent.WithName("Subordinate"))

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithName("TestAgent"),
		agent.WithSubordinate(subordinate),
		agent.WithTool(tool),
	)

	_, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)
	assert.NoError(t, err)
}

func TestAgent_Run_ErrorEmptyFinishReason(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			return types.Response{}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(model.Wrap())

	_, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)
	require.Error(t, err)
	assert.Equal(t, agent.ErrModelResponse{Msg: "no finish reason"}, err)
}

func TestAgent_Run_ContinueUntilNoToolsCalled(t *testing.T) {
	chatCalls := 0
	funcCalls := 0

	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			if chatCalls < 3 {
				chatCalls++
				return types.Response{
					FinishReason: "tool_call",
					ToolCalls: []types.ToolCall{
						{
							ID:        strconv.Itoa(chatCalls),
							Name:      "test_tool",
							Arguments: "",
						},
					},
				}, nil
			}
			return types.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	tool := tools.New("test_tool", "A test tool", func(ctx context.Context) ([]types.Content, error) {
		funcCalls++
		content := types.NewTextContent("tool response " + strconv.Itoa(funcCalls))
		return content.AsSlice(), nil
	})

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithTool(tool),
	)

	_, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)
	require.NoError(t, err)
	assert.Equal(t, 3, chatCalls)
	assert.Equal(t, 3, funcCalls)
}

func TestAgent_Run_AssistantMessageAppended(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			return types.Response{
				FinishReason: "stop",
				Content:      types.NewTextContent("This is the assistant response."),
			}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(model.Wrap())

	run, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)

	require.NoError(t, err)
	require.NotNil(t, run)

	require.Equal(t, 2, len(run.Messages))
	assert.Equal(t, types.Message{Role: types.RoleUser, Content: types.NewTextContent("Hello")}, run.Messages[0])
	assert.Equal(t, types.Message{Role: types.RoleAssistant, Content: types.NewTextContent("This is the assistant response.")}, run.Messages[1])
}

func TestAgent_Run_MetricsPopulated(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			return types.Response{
				FinishReason: "stop",
				Content:      types.NewTextContent("This is the assistant response."),
				Metrics: metrics.Metrics{
					PromptTokens:     10,
					CompletionTokens: 20,
					TotalTokens:      30,
				},
			}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithID(uuid.Max),
	)

	run, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)

	require.NoError(t, err)
	require.NotNil(t, run)
	require.NotNil(t, run.Metrics)

	metricsList, exists := run.Metrics[uuid.Max]
	require.True(t, exists)
	require.Equal(t, 1, len(metricsList))
	assert.Equal(t, uint(10), metricsList[0].PromptTokens)
	assert.Equal(t, uint(20), metricsList[0].CompletionTokens)
	assert.Equal(t, uint(30), metricsList[0].TotalTokens)
}

func TestAgent_Run_ToolErrorHandling(t *testing.T) {
	toolCalled := false
	model := MockModel{
		ChatFunc: func(ctx context.Context, messages []types.Message, tools []tools.Tool, opts *types.Options) (types.Response, error) {
			if toolCalled {
				return types.Response{
					FinishReason: "stop",
					Content:      types.NewTextContent("Finished after tool call."),
				}, nil
			}
			toolCalled = true
			return types.Response{
				FinishReason: "tool_call",
				ToolCalls: []types.ToolCall{
					{
						ID:        "call_id_1",
						Name:      "error_tool",
						Arguments: "",
					},
				},
			}, nil
		},
		ChatStreamingFunc: nil,
	}

	tool := tools.New("error_tool", "A tool that returns an error", func(ctx context.Context) ([]types.Content, error) {
		return nil, assert.AnError
	})

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithTool(tool),
	)

	run, err := agentInstance.RunWithParams(
		t.Context(),
		types.WithTextMessage(types.RoleUser, "Hello"),
	)
	require.NoError(t, err)
	require.NotNil(t, run)
	require.Equal(t, 4, len(run.Messages))
	assert.Equal(t, types.Message{Role: types.RoleUser, Content: types.NewTextContent("Hello")}, run.Messages[0])
	assert.Equal(t, types.Message{Role: types.RoleAssistant, Content: types.NewTextContent(""), ToolCalls: []types.ToolCall{{ID: "call_id_1", Name: "error_tool", Arguments: ""}}}, run.Messages[1])
	assert.Equal(t, types.Message{Role: types.RoleTool, Name: "call_id_1", Content: types.NewTextContent("error: " + assert.AnError.Error()), IsErr: true}, run.Messages[2])
	assert.Equal(t, types.Message{Role: types.RoleAssistant, Content: types.NewTextContent("Finished after tool call.")}, run.Messages[3])
}
