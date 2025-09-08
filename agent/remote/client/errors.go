package client

type ErrUnknownEvent struct {
	Event string
}

func (e ErrUnknownEvent) Error() string {
	return "unknown event type: " + e.Event
}

type ErrChunkMessage struct {
	Message string
}

func (e ErrChunkMessage) Error() string {
	return "chunk message error: " + e.Message
}

type ErrServerResponse struct {
	Err        error
	StatusCode int
}

func (e ErrServerResponse) Error() string {
	return "client request error: " + e.Err.Error()
}
