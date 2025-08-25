package openai

import (
	"time"

	"github.com/openai/openai-go"
	"github.com/yourlogarithm/l337/metrics"
)

func convertMetrics(openaiMetrics *openai.CompletionUsage, totalDuration time.Duration) metrics.Metrics {
	if openaiMetrics == nil {
		return metrics.Metrics{}
	}
	return metrics.Metrics{
		Timestamp:     time.Now(),
		TotalDuration: totalDuration,
		TotalTokens:   uint(openaiMetrics.TotalTokens),

		PromptAudioTokens:  uint(openaiMetrics.PromptTokensDetails.AudioTokens),
		PromptCachedTokens: uint(openaiMetrics.PromptTokensDetails.CachedTokens),
		PromptTokens:       uint(openaiMetrics.PromptTokens),

		CompletionAcceptedPredictionTokens: uint(openaiMetrics.CompletionTokensDetails.AcceptedPredictionTokens),
		CompletionAudioTokens:              uint(openaiMetrics.CompletionTokensDetails.AudioTokens),
		CompletionReasoningTokens:          uint(openaiMetrics.CompletionTokensDetails.ReasoningTokens),
		CompletionRejectedPredictionTokens: uint(openaiMetrics.CompletionTokensDetails.RejectedPredictionTokens),
		CompletionTokens:                   uint(openaiMetrics.CompletionTokens),
	}
}
