package provider

import "github.com/yourlogarithm/l337/chat"

type ContentAccumulator struct {
	Response        Response
	toolCallsLookup map[string]*chat.ToolCall
}

func (acc *ContentAccumulator) AddChunk(chunk *Response) error {
	if acc.Response.ID == "" {
		acc.Response.ID = chunk.ID
	} else if acc.Response.ID != chunk.ID && chunk.ID != "" {
		return ErrChunkAddition{Accumulator: acc.Response, Chunk: chunk}
	}

	acc.Response.Created = chunk.Created
	acc.Response.Content.Text += chunk.Content.Text
	acc.Response.Refusal += chunk.Refusal
	acc.Response.Reasoning += chunk.Reasoning
	acc.Response.FinishReason += chunk.FinishReason
	acc.Response.Metrics.Add(&chunk.Metrics)

	for _, call := range chunk.ToolCalls {
		if accToolCall, exists := acc.toolCallsLookup[call.ID]; exists {
			accToolCall.Arguments += call.Arguments
		} else if call.ID == "" {
			if len(acc.Response.ToolCalls) == 0 {
				return ErrChunkAddition{Accumulator: acc.Response, Chunk: chunk}
			}
			lastCall := &acc.Response.ToolCalls[len(acc.Response.ToolCalls)-1]
			lastCall.Arguments += call.Arguments
		} else {
			idx := len(acc.Response.ToolCalls)
			acc.Response.ToolCalls = append(acc.Response.ToolCalls, call)
			acc.toolCallsLookup[call.ID] = &acc.Response.ToolCalls[idx]
		}
	}

	return nil
}

func NewContentAccumulator() ContentAccumulator {
	return ContentAccumulator{
		toolCallsLookup: make(map[string]*chat.ToolCall),
	}
}

func (acc ContentAccumulator) HasFinished() bool {
	return acc.Response.FinishReason != ""
}
