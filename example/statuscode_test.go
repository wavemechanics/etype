package example

import (
	"fmt"
	"testing"

	"github.com/wavemechanics/etype"
	"golang.org/x/xerrors"
)

func TestStatusCode(t *testing.T) {

	// add a status code to some error
	const cause = etype.Sentinel("original error")
	err := WithStatusCode(cause, 400)

	// tests both unwrap and Is
	if !xerrors.Is(err, cause) {
		t.Errorf("StatusCode Is: expected true")
	}

	// make sure As works
	err = fmt.Errorf("generic annotation: %w", err)
	var errStatusCode *ErrStatusCode
	if !xerrors.As(err, &errStatusCode) {
		t.Errorf("StatusCode: As didn't work")
	}

	// now that we've found the ErrStatusCode...
	if errStatusCode.StatusCode() != 400 {
		t.Errorf("StatusCode didn't return 400")
	}
}
