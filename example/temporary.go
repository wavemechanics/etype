package example

import (
	"fmt"

	"golang.org/x/xerrors"
)

// A TemporaryError is an example of a custom error that can be tested
// with .Is.
type TemporaryError struct {
	error
}

// NewTemporary wraps an existing error, giving it type TemporaryError.
func NewTemporary(err error) error {
	return &TemporaryError{fmt.Errorf("temporary: %w", err)}
}

// Is returns true if err is of type NewTemporary.
func (e *TemporaryError) Is(err error) bool {
	_, ok := err.(*TemporaryError)
	return ok
}

// Unwrap returns the underlying error.
func (e *TemporaryError) Unwrap() error {
	return xerrors.Unwrap(e.error)
}

// IsTemporary returns true if any error in the chain is a TemporaryError.
// This is just a shortcut to using the more verbose xerrors.Is call.
func IsTemporary(err error) bool {
	return xerrors.Is(err, &TemporaryError{})
}
