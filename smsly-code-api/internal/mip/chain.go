package mip

import "github.com/google/uuid"

// ChainVerification represents the result of verifying a chain of stamps.
type ChainVerification struct {
	Valid    bool       `json:"valid"`
	BrokenAt *uuid.UUID `json:"broken_at,omitempty"`
	Error    string     `json:"error,omitempty"`
}
