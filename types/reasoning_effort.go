package types

type ReasoningEffortLevel string

const (
	ReasoningEffortLow    ReasoningEffortLevel = "low"
	ReasoningEffortMedium ReasoningEffortLevel = "medium"
	ReasoningEffortHigh   ReasoningEffortLevel = "high"
)

type ReasoningEffortUnion struct {
	ofString ReasoningEffortLevel
	ofBool   *bool
}

func NewReasoningEffortLevel(level ReasoningEffortLevel) *ReasoningEffortUnion {
	return &ReasoningEffortUnion{ofString: level}
}

func NewReasoningEffortBool(enabled bool) *ReasoningEffortUnion {
	return &ReasoningEffortUnion{ofBool: &enabled}
}

func (r *ReasoningEffortUnion) AsLevel() (ReasoningEffortLevel, bool) {
	if r.ofString == "" {
		return "", false
	}
	return r.ofString, true
}

func (r *ReasoningEffortUnion) AsBool() (bool, bool) {
	if r.ofBool == nil {
		return false, false
	}
	return *r.ofBool, true
}

func (r *ReasoningEffortUnion) AsAny() any {
	if level, ok := r.AsLevel(); ok {
		return level
	}
	if enabled, ok := r.AsBool(); ok {
		return enabled
	}
	return nil
}
