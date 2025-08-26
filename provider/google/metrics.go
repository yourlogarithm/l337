package google

import (
	"time"

	"github.com/yourlogarithm/l337/metrics"
	"google.golang.org/genai"
)

func convertMetrics(googleMetrics *genai.GenerateContentResponseUsageMetadata, totalDuration time.Duration) metrics.Metrics {
	return metrics.Metrics{
		Timestamp: time.Now(),

		TotalDuration: totalDuration,
		TotalTokens:   uint(googleMetrics.TotalTokenCount),

		CompletionTokens: uint(googleMetrics.TotalTokenCount - googleMetrics.PromptTokenCount),

		PromptTokens:       uint(googleMetrics.PromptTokenCount),
		PromptCachedTokens: uint(googleMetrics.CachedContentTokenCount),
	}
}
