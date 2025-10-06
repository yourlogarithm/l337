package agent

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

func BuildRunResponse(params ...chat.Parameter) (*chat.RunResponse, error) {
	runResponse := &chat.RunResponse{
		Metrics: make(map[uuid.UUID][]metrics.Metrics),
	}
	for _, param := range params {
		if err := param.Apply(runResponse); err != nil {
			return nil, err
		}
	}
	return runResponse, nil
}
