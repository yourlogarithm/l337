package streaming

import (
	"io"

	"github.com/yourlogarithm/l337/types"
)

type ResponseChannel chan Chunk

func NewResponseChannel(size ...int) ResponseChannel {
	var bufSize int
	if len(size) > 0 {
		bufSize = size[0]
	}
	return make(chan Chunk, bufSize)
}

func (rs ResponseChannel) Next() (*types.Response, error) {
	p, ok := <-rs
	if !ok {
		return nil, io.EOF
	}
	return p.Chunk, p.Error
}

func (rs ResponseChannel) Send(data *types.Response) {
	rs <- Chunk{Chunk: data}
}

func (rs ResponseChannel) SendErr(err error) {
	rs <- Chunk{Error: err}
}

func (rs ResponseChannel) Close() {
	close(rs)
}

func (rs ResponseChannel) Drain() int {
	count := 0
	for range rs {
		count++
	}
	return count
}
