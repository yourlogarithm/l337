package provider

import (
	"fmt"
)

type ErrUnknownRole struct {
	Role string
}

func (e ErrUnknownRole) Error() string {
	return fmt.Sprintf("unknown role: %s", e.Role)
}

type ErrChunkAddition struct {
	Accumulator Response
	Chunk       *Response
}

func (e ErrChunkAddition) Error() string {
	return fmt.Sprintf("failed to add chunk %v to accumulator %v", e.Chunk, e.Accumulator)
}
