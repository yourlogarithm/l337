package tools

type SkillCard struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Examples    []string `json:"examples"`
}

type SkillCardOption interface {
	Apply(*SkillCard)
}

type SkillOptionFunc func(*SkillCard)

func (s SkillOptionFunc) Apply(r *SkillCard) { s(r) }

func WithTags(tags ...string) SkillCardOption {
	return SkillOptionFunc(func(o *SkillCard) {
		o.Tags = append(o.Tags, tags...)
	})
}

func WithExamples(examples ...string) SkillCardOption {
	return SkillOptionFunc(func(o *SkillCard) {
		o.Examples = append(o.Examples, examples...)
	})
}
