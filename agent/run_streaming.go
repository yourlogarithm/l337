package agent

import (
	"context"
	"io"
	"slices"

	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/provider"
	"github.com/yourlogarithm/l337/tools"
)

func (a *Agent) RunStreamingWithParams(ctx context.Context, bufferSize int, params ...chat.Parameter) (*chat.RunResponse, provider.ResponseChannel, error) {
	runResponse, err := BuildRunResponse(params...)
	if err != nil {
		return nil, nil, err
	}
	channel, err := a.RunStreaming(ctx, runResponse, bufferSize)
	return runResponse, channel, err
}

func (a *Agent) RunStreaming(ctx context.Context, runResponse *chat.RunResponse, bufferSize int) (provider.ResponseChannel, error) {
	if len(runResponse.Messages) == 0 || runResponse.Messages[0].Role != chat.RoleSystem {
		content, err := a.ComputeSystemMessage()
		if err != nil {
			return nil, err
		}
		systemMsg := chat.Message{
			Role:    chat.RoleSystem,
			Content: content,
		}
		runResponse.Messages = slices.Insert(runResponse.Messages, 0, systemMsg)
	}

	channel := provider.NewResponseChannel(bufferSize)

	go a.scheduleChunkProcessing(ctx, runResponse, channel)

	return channel, nil
}

func (a *Agent) scheduleChunkProcessing(
	ctx context.Context,
	runResponse *chat.RunResponse,
	responseChannel provider.ResponseChannel,
) {
	defer responseChannel.Close()

	tools := make([]tools.Tool, 0, len(a.tools))
	for _, tool := range a.tools {
		tools = append(tools, tool)
	}

	for {
		req := provider.Request{
			Messages: runResponse.Messages,
			Tools:    tools,
		}
		logger.Debug("agent.run.request", "agent", a.name, "request", req)
		stream, err := a.model.Impl.ChatStreaming(ctx, &req, &a.chatOptions)
		if err != nil {
			responseChannel.SendErr(err)
			return
		}

		acc := provider.ContentAccumulator{}
		for {
			chunk, err := stream.Next()
			if err == io.EOF {
				break
			} else if err != nil {
				responseChannel.SendErr(err)
				return
			}
			if err = acc.AddChunk(chunk); err != nil {
				responseChannel.SendErr(err)
				return
			}
			responseChannel.Send(chunk)
		}

		toolsCalled, err := a.handleResponse(ctx, runResponse, &acc.Response)
		if err != nil {
			responseChannel.SendErr(err)
			return
		}

		if !toolsCalled {
			break
		}
	}
}
