package anthropic

import (
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/yourlogarithm/l337/metrics"
)

func convertMetrics(anthropicMetrics *anthropic.Usage, totalDuration time.Duration) metrics.Metrics {
	return metrics.Metrics{
		Timestamp:     time.Now(),
		TotalDuration: totalDuration,
		TotalTokens:   uint(anthropicMetrics.InputTokens + anthropicMetrics.OutputTokens),

		CompletionTokens: uint(anthropicMetrics.OutputTokens),

		PromptCacheCreationTokens: uint(anthropicMetrics.CacheCreationInputTokens),
		PromptCachedTokens:        uint(anthropicMetrics.CacheReadInputTokens),
		PromptTokens:              uint(anthropicMetrics.InputTokens),
	}
}
