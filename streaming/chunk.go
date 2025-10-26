package streaming

import "github.com/yourlogarithm/l337/types"

type Chunk struct {
	Chunk *types.Response
	Error error
}

type ChunkMarshalable struct {
	Chunk *types.Response `json:"chunk,omitempty"`
	Error string          `json:"error,omitempty"`
}

func (rc *Chunk) ToMarshalable() ChunkMarshalable {
	if rc.Error == nil {
		return ChunkMarshalable{
			Chunk: rc.Chunk,
		}
	}
	return ChunkMarshalable{
		Chunk: rc.Chunk,
		Error: rc.Error.Error(),
	}
}
