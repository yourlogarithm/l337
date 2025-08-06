package team

import "fmt"

type TeamError string

const (
	ErrMemberNotFound TeamError = "member %s not found"
)

func (e TeamError) Error(args ...any) error {
	return fmt.Errorf(string(e), args...)
}
