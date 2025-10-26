package streaming

import (
	"fmt"

	"github.com/yourlogarithm/l337/types"
)

type ErrChunkAddition struct {
	Accumulator types.Response
	Chunk       *types.Response
}

func (e ErrChunkAddition) Error() string {
	return fmt.Sprintf("failed to add chunk %v to accumulator %v", e.Chunk, e.Accumulator)
}
