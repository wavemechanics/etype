package example

import (
	"golang.org/x/xerrors"
)

// MyError is an example of a custom error type that can be tested with .Is.
type MyError struct {
	error
}

// NewMyError wraps an existing error, giving it type MyError
func NewMyError(err error) error {
	return &MyError{xerrors.Errorf("my error: %w", err)}
}

// Unwrap returns the underlying error.
func (e *MyError) Unwrap() error {
	return xerrors.Unwrap(e.error)
}

// Is returns true if err is of type MyError.
func (e *MyError) Is(err error) bool {
	_, ok := err.(*MyError)
	return ok
}
