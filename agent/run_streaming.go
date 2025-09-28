package agent

import (
	"context"
	"io"

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
	tools, err := a.preRun(runResponse)
	if err != nil {
		return nil, err
	}

	channel := provider.NewResponseChannel(bufferSize)

	go a.scheduleChunkProcessing(ctx, tools, runResponse, channel)

	return channel, nil
}

func (a *Agent) scheduleChunkProcessing(
	ctx context.Context,
	tools []tools.Tool,
	runResponse *chat.RunResponse,
	responseChannel provider.ResponseChannel,
) {
	defer responseChannel.Close()

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

		acc := provider.NewContentAccumulator()
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
