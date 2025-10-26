package agent

import (
	"github.com/google/uuid"
	"github.com/yourlogarithm/l337/metrics"
	"github.com/yourlogarithm/l337/types"
)

func BuildRun(params ...types.Parameter) (*types.Run, error) {
	run := &types.Run{
		Metrics: make(map[uuid.UUID][]metrics.Metrics),
	}
	for _, param := range params {
		if err := param.Apply(run); err != nil {
			return nil, err
		}
	}
	return run, nil
}
