package agent_test

import (
	"context"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yourlogarithm/l337/agent"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

func TestAgent_Run_NilModel(t *testing.T) {
	agentInstance, err := agent.New(nil, agent.WithName("TestAgent"))
	require.NoError(t, err)

	_, err = agentInstance.RunWithParams(t.Context(), chat.WithTextMessage(chat.RoleUser, "Hello"))
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
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			assert.Equal(t, 1, len(req.Messages))
			assert.Equal(t, chat.Message{Role: chat.RoleUser, Content: chat.NewTextContent("Hello")}, req.Messages[0])
			return provider.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(model.Wrap())

	_, err := agentInstance.RunWithParams(
		t.Context(),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)

	assert.NoError(t, err)
}

func TestAgent_Run_SystemMessageAppendedIfNotProvided(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			assert.Equal(t, 2, len(req.Messages))
			assert.Equal(t, chat.Message{Role: chat.RoleSystem, Content: chat.NewTextContent("Your name is TestAgent.\nA test agent")}, req.Messages[0])
			assert.Equal(t, chat.Message{Role: chat.RoleUser, Content: chat.NewTextContent("Hello")}, req.Messages[1])
			return provider.Response{FinishReason: "stop"}, nil
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
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)
	assert.NoError(t, err)
}

func TestAgent_Run_SystemMessageNotAppendedIfProvided(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			assert.Equal(t, 2, len(req.Messages))
			assert.Equal(t, chat.Message{Role: chat.RoleSystem, Content: chat.NewTextContent("Custom system message")}, req.Messages[0])
			assert.Equal(t, chat.Message{Role: chat.RoleUser, Content: chat.NewTextContent("Hello")}, req.Messages[1])
			return provider.Response{FinishReason: "stop"}, nil
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
		chat.WithTextMessage(chat.RoleSystem, "Custom system message"),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)
	assert.NoError(t, err)
}

func TestAgent_Run_ToolsPassed(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			assert.Equal(t, 2, len(req.Tools)) // Including the delegate_task tool

			toolMap := make(map[string]tools.Tool)
			for _, t := range req.Tools {
				toolMap[t.Name] = t
			}

			assert.Equal(t, "test_tool", toolMap["test_tool"].Name)
			assert.Equal(t, "A test tool", toolMap["test_tool"].Description)

			assert.Equal(t, agent.DELEGATE_TASK_TOOL_NAME, toolMap[agent.DELEGATE_TASK_TOOL_NAME].Name)
			assert.Equal(t, agent.DELEGATE_TASK_TOOL_DESC, toolMap[agent.DELEGATE_TASK_TOOL_NAME].Description)

			return provider.Response{FinishReason: "stop"}, nil
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
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)
	assert.NoError(t, err)
}

func TestAgent_Run_ErrorEmptyFinishReason(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			return provider.Response{}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(model.Wrap())

	_, err := agentInstance.RunWithParams(
		t.Context(),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)
	require.Error(t, err)
	assert.Equal(t, agent.ErrModelResponse{Msg: "no finish reason"}, err)
}

func TestAgent_Run_ContinueUntilNoToolsCalled(t *testing.T) {
	chatCalls := 0
	funcCalls := 0

	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			if chatCalls < 3 {
				chatCalls++
				return provider.Response{
					FinishReason: "tool_call",
					ToolCalls: []chat.ToolCall{
						{
							ID:        strconv.Itoa(chatCalls),
							Name:      "test_tool",
							Arguments: "",
						},
					},
				}, nil
			}
			return provider.Response{FinishReason: "stop"}, nil
		},
		ChatStreamingFunc: nil,
	}

	tool := tools.New("test_tool", "A test tool", func(ctx context.Context) ([]chat.Content, error) {
		funcCalls++
		content := chat.NewTextContent("tool response " + strconv.Itoa(funcCalls))
		return content.AsSlice(), nil
	})

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithTool(tool),
	)

	_, err := agentInstance.RunWithParams(
		t.Context(),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)
	require.NoError(t, err)
	assert.Equal(t, 3, chatCalls)
	assert.Equal(t, 3, funcCalls)
}

func TestAgent_Run_AssistantMessageAppended(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			return provider.Response{
				FinishReason: "stop",
				Content:      "This is the assistant response.",
			}, nil
		},
		ChatStreamingFunc: nil,
	}

	agentInstance, _ := agent.New(model.Wrap())

	runResponse, err := agentInstance.RunWithParams(
		t.Context(),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)

	require.NoError(t, err)
	require.NotNil(t, runResponse)

	require.Equal(t, 2, len(runResponse.Messages))
	assert.Equal(t, chat.Message{Role: chat.RoleUser, Content: chat.NewTextContent("Hello")}, runResponse.Messages[0])
	assert.Equal(t, chat.Message{Role: chat.RoleAssistant, Content: chat.NewTextContent("This is the assistant response.")}, runResponse.Messages[1])
}

func TestAgent_Run_MetricsPopulated(t *testing.T) {
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			return provider.Response{
				FinishReason: "stop",
				Content:      "This is the assistant response.",
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

	runResponse, err := agentInstance.RunWithParams(
		t.Context(),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)

	require.NoError(t, err)
	require.NotNil(t, runResponse)
	require.NotNil(t, runResponse.Metrics)

	metricsList, exists := runResponse.Metrics[uuid.Max]
	require.True(t, exists)
	require.Equal(t, 1, len(metricsList))
	assert.Equal(t, uint(10), metricsList[0].PromptTokens)
	assert.Equal(t, uint(20), metricsList[0].CompletionTokens)
	assert.Equal(t, uint(30), metricsList[0].TotalTokens)
}

func TestAgent_Run_ToolErrorHandling(t *testing.T) {
	toolCalled := false
	model := MockModel{
		ChatFunc: func(ctx context.Context, req *provider.Request, opts *chat.Options) (provider.Response, error) {
			if toolCalled {
				return provider.Response{
					FinishReason: "stop",
					Content:      "Finished after tool call.",
				}, nil
			}
			toolCalled = true
			return provider.Response{
				FinishReason: "tool_call",
				ToolCalls: []chat.ToolCall{
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

	tool := tools.New("error_tool", "A tool that returns an error", func(ctx context.Context) ([]chat.Content, error) {
		return nil, assert.AnError
	})

	agentInstance, _ := agent.New(
		model.Wrap(),
		agent.WithTool(tool),
	)

	runResponse, err := agentInstance.RunWithParams(
		t.Context(),
		chat.WithTextMessage(chat.RoleUser, "Hello"),
	)
	require.NoError(t, err)
	require.NotNil(t, runResponse)
	require.Equal(t, 4, len(runResponse.Messages))
	assert.Equal(t, chat.Message{Role: chat.RoleUser, Content: chat.NewTextContent("Hello")}, runResponse.Messages[0])
	assert.Equal(t, chat.Message{Role: chat.RoleAssistant, Content: chat.NewTextContent(""), ToolCalls: []chat.ToolCall{{ID: "call_id_1", Name: "error_tool", Arguments: ""}}}, runResponse.Messages[1])
	assert.Equal(t, chat.Message{Role: chat.RoleTool, Name: "call_id_1", Content: chat.NewTextContent("error: " + assert.AnError.Error()), IsErr: true}, runResponse.Messages[2])
	assert.Equal(t, chat.Message{Role: chat.RoleAssistant, Content: chat.NewTextContent("Finished after tool call.")}, runResponse.Messages[3])
}
