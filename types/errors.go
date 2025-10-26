package types

import (
	"fmt"
)

// --- ErrUnknownRole ---

type ErrUnknownRole struct {
	Role string
}

func (e ErrUnknownRole) Error() string {
	return fmt.Sprintf("unknown role: %s", e.Role)
}

// --- ErrParams ---
