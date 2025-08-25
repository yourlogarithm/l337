package ollama

import (
	"time"

	"github.com/ollama/ollama/api"
	"github.com/yourlogarithm/l337/metrics"
)

func convertMetrics(ollamaMetrics *api.Metrics) metrics.Metrics {
	if ollamaMetrics == nil {
		return metrics.Metrics{}
	}
	return metrics.Metrics{
		Timestamp:     time.Now(),
		LoadDuration:  ollamaMetrics.LoadDuration,
		TotalDuration: ollamaMetrics.TotalDuration,
		TotalTokens:   uint(ollamaMetrics.PromptEvalCount + ollamaMetrics.EvalCount),

		CompletionTokens: uint(ollamaMetrics.EvalCount),
		PromptTokens:     uint(ollamaMetrics.PromptEvalCount),
	}
}
