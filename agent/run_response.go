package agent

import (
	"fmt"

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
	if len(runParams.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	runResponse := chat.RunResponse{
		SessionID: runParams.SessionID,
		Messages:  runParams.Messages,
		Metrics:   make(map[uuid.UUID][]metrics.Metrics),
	}

	return &runResponse, nil
}
