package tools

type SkillCard struct {
	Name        string
	Description string
	Tags        []string
	Examples    []string
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
