package provider

import "github.com/yourlogarithm/l337/chat"

type ContentAccumulator struct {
	Response
}

func (acc *ContentAccumulator) AddChunk(chunk *Response) error {
	if acc.ID == "" {
		acc.ID = chunk.ID
	} else if acc.ID != chunk.ID && chunk.ID != "" {
		return ErrChunkAddition{Accumulator: acc.Response, Chunk: chunk}
	}

	acc.Created = chunk.Created
	acc.Content += chunk.Content
	acc.Refusal += chunk.Refusal
	acc.Reasoning += chunk.Reasoning
	acc.FinishReason += chunk.FinishReason
	acc.Metrics.Add(&chunk.Metrics)

	if len(acc.ToolCalls) < len(chunk.ToolCalls) {
		acc.ToolCalls = make([]chat.ToolCall, len(chunk.ToolCalls))
	}

	for i := range chunk.ToolCalls {
		call := &chunk.ToolCalls[i]
		accCall := &acc.ToolCalls[i]
		if accCall.ID == "" {
			*accCall = *call
		} else if accCall.ID != call.ID && call.ID != "" {
			return &ErrChunkAddition{Accumulator: acc.Response, Chunk: chunk}
		}
		if call.Name != "" {
			accCall.Name = call.Name
		}
		accCall.Arguments += call.Arguments
	}

	return nil
}

func (acc *ContentAccumulator) HasFinished() bool {
	return acc.Response.FinishReason != ""
}
