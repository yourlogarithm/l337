package agent

import (
	"context"
	"io"

	"github.com/yourlogarithm/l337/streaming"
	"github.com/yourlogarithm/l337/tools"
	"github.com/yourlogarithm/l337/types"
)

func (a *Agent) RunStreamingWithParams(ctx context.Context, bufferSize int, params ...types.Parameter) (*types.Run, streaming.ResponseChannel, error) {
	run, err := BuildRun(params...)
	if err != nil {
		return nil, nil, err
	}
	channel, err := a.RunStreaming(ctx, run, bufferSize)
	return run, channel, err
}

func (a *Agent) RunStreaming(ctx context.Context, run *types.Run, bufferSize int) (streaming.ResponseChannel, error) {
	tools, err := a.preRun(run)
	if err != nil {
		return nil, err
	}

	channel := streaming.NewResponseChannel(bufferSize)

	go a.scheduleChunkProcessing(ctx, tools, run, channel)

	return channel, nil
}

func (a *Agent) scheduleChunkProcessing(
	ctx context.Context,
	tools []tools.Tool,
	run *types.Run,
	responseChannel streaming.ResponseChannel,
) {
	defer responseChannel.Close()

	for {
		logger.Debug("agent.run.request", "agent", a.name, "messages", run.Messages, "tools", tools)
		stream, err := a.model.Impl.ChatStreaming(ctx, run.Messages, tools, &a.chatOptions)
		if err != nil {
			responseChannel.SendErr(err)
			return
		}

		acc := streaming.NewContentAccumulator()
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

		toolsCalled, err := a.handleResponse(ctx, run, &acc.Response)
		if err != nil {
			responseChannel.SendErr(err)
			return
		}

		if !toolsCalled {
			break
		}
	}
}
