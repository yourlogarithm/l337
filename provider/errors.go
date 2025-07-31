package provider

import (
	"errors"
	"fmt"
)

// ErrUnknownRole is used as a sentinel error for unknown roles.
var ErrUnknownRole = errors.New("unknown role")

// NewUnknownRoleError creates an error for an unknown role that wraps ErrUnknownRole.
func NewUnknownRoleError(role string) error {
	return fmt.Errorf("%w: %s", ErrUnknownRole, role)
}
