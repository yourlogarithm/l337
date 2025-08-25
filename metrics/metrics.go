package metrics

import (
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	Timestamp time.Time `json:"timestamp"`
	SessionID uuid.UUID `json:"session_id"`

	LoadDuration  time.Duration `json:"load_duration"`
	TotalDuration time.Duration `json:"total_duration"`
	TotalTokens   uint          `json:"total_tokens"`

	CompletionAcceptedPredictionTokens uint `json:"completion_accepted_prediction_tokens"`
	CompletionAudioTokens              uint `json:"completion_audio_tokens"`
	CompletionReasoningTokens          uint `json:"completion_reasoning_tokens"`
	CompletionRejectedPredictionTokens uint `json:"completion_rejected_prediction_tokens"`
	CompletionTokens                   uint `json:"completion_tokens"`

	PromptAudioTokens  uint `json:"prompt_audio_tokens"`
	PromptCachedTokens uint `json:"prompt_cached_tokens"`
	PromptTokens       uint `json:"prompt_tokens"`
	// Anthropic [Usage.CacheCreationInputTokens](https://pkg.go.dev/github.com/anthropics/anthropic-sdk-go@v1.7.0#Usage.CacheCreationInputTokens)
	PromptCacheCreationTokens uint `json:"prompt_cache_creation_tokens"`
}

func (m *Metrics) Add(other *Metrics) {
	if other == nil {
		return
	}
	m.TotalTokens += other.TotalTokens
	m.TotalDuration += other.TotalDuration
	m.LoadDuration += other.LoadDuration

	m.CompletionTokens += other.CompletionTokens
	m.CompletionAcceptedPredictionTokens += other.CompletionAcceptedPredictionTokens
	m.CompletionAudioTokens += other.CompletionAudioTokens
	m.CompletionReasoningTokens += other.CompletionReasoningTokens
	m.CompletionRejectedPredictionTokens += other.CompletionRejectedPredictionTokens

	m.PromptTokens += other.PromptTokens
	m.PromptAudioTokens += other.PromptAudioTokens
	m.PromptCachedTokens += other.PromptCachedTokens
}

func (m *Metrics) TokensPerSecond() float64 {
	if m.TotalDuration == 0 {
		return 0
	}
	return float64(m.TotalTokens) / m.TotalDuration.Seconds()
}
