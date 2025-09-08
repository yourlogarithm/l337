package provider

import "io"

type ResponseChunkMarshalable struct {
	Chunk *Response `json:"chunk,omitempty"`
	Error string    `json:"error,omitempty"`
}

type ResponseChunk struct {
	Chunk *Response
	Error error
}

func (rc *ResponseChunk) ToMarshalable() ResponseChunkMarshalable {
	if rc.Error == nil {
		return ResponseChunkMarshalable{
			Chunk: rc.Chunk,
		}
	}
	return ResponseChunkMarshalable{
		Chunk: rc.Chunk,
		Error: rc.Error.Error(),
	}
}

type ResponseChannel chan ResponseChunk

func NewResponseChannel(size ...int) ResponseChannel {
	var bufSize int
	if len(size) > 0 {
		bufSize = size[0]
	}
	return make(chan ResponseChunk, bufSize)
}

func (rs ResponseChannel) Next() (*Response, error) {
	p, ok := <-rs
	if !ok {
		return nil, io.EOF
	}
	return p.Chunk, p.Error
}

func (rs ResponseChannel) Send(data *Response) {
	rs <- ResponseChunk{Chunk: data}
}

func (rs ResponseChannel) SendErr(err error) {
	rs <- ResponseChunk{Error: err}
}

func (rs ResponseChannel) Close() {
	close(rs)
}
