package example

import (
	"errors"
	"fmt"
)

// ErrStatusCode is an error that holds some kind of status code.
type ErrStatusCode struct {
	error
	statusCode int
}

// WithStatusCode wraps an existing error, giving it a status code.
func WithStatusCode(err error, code int) error {
	return &ErrStatusCode{
		error:      fmt.Errorf("code %d: %w", code, err),
		statusCode: code,
	}
}

// Unwrap returns the underlying error.
func (e *ErrStatusCode) Unwrap() error {
	return errors.Unwrap(e.error)
}

// Is is true if err is of type ErrStatysCode.
func (e *ErrStatusCode) Is(err error) bool {
	_, ok := err.(*ErrStatusCode)
	return ok
}

// StatusCode returns the status code associated with this error.
// If .As returns true that an error is an ErrStatusCode, then this
// method may be used on the result.
func (e *ErrStatusCode) StatusCode() int {
	return e.statusCode
}
