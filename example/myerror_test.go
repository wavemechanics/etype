package example

import (
	"testing"

	"github.com/wavemechanics/etype"
	"golang.org/x/xerrors"
)

func TestMyError(t *testing.T) {

	const cause = etype.Sentinel("original error")
	err := NewMyError(cause)

	// tests both Unwrap and Is
	if !xerrors.Is(err, cause) {
		t.Errorf("MyError Is: expected true")
	}
}
