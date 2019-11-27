package example

import (
	"golang.org/x/xerrors"
	"testing"

	"github.com/wavemechanics/etype"
)

func TestTemporary(t *testing.T) {

	// wrap some error in a TemporaryError
	const cause = etype.Sentinel("original cause")
	err := NewTemporary(cause)

	// don't expect the TemporaryError to obscure the original error
	if !xerrors.Is(err, cause) {
		t.Errorf("Is: expected to match original cause")
	}

	// annotate, and still expect it to be a TemporaryError
	err = xerrors.Errorf("annotated: %w", err)
	if !IsTemporary(err) {
		t.Errorf("IsTemporary: expected annoted error to also be Temporary")
	}
}
