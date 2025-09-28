package agent

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/chat"
	"github.com/yourlogarithm/l337/metrics"
)

func BuildRunResponse(params ...chat.Parameter) (*chat.RunResponse, error) {
	var runParams chat.Parameters
	for _, param := range params {
		if err := param.Apply(&runParams); err != nil {
			return nil, err
		}
	}

	runResponse := chat.RunResponse{
		SessionID: runParams.SessionID,
		Messages:  runParams.Messages,
		Metrics:   make(map[uuid.UUID][]metrics.Metrics),
	}

	return &runResponse, nil
}
