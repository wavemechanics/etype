package etype_test

import (
	"testing"

	"github.com/wavemechanics/etype"

	"golang.org/x/xerrors"
)

func TestSentinel(t *testing.T) {
	const errExample = etype.Sentinel("an example error")

	var tests = []struct {
		desc       string
		err        error
		isSentinel bool
	}{
		{
			desc:       "nil is not a Sentinel",
			err:        nil,
			isSentinel: false,
		},
		{
			desc:       "non-Sentinel is not a sentinel",
			err:        xerrors.New("an example error"),
			isSentinel: false,
		},
		{
			desc:       "Sentinel of a different value does not match",
			err:        etype.Sentinel("a different error"),
			isSentinel: false,
		},
		{
			desc:       "Sentinel of the same value matches",
			err:        etype.Sentinel("an example error"),
			isSentinel: true,
		},
		{
			desc:       "The same Sentinel matches",
			err:        errExample,
			isSentinel: true,
		},
	}

	for _, test := range tests {
		got := xerrors.Is(test.err, errExample)
		if got != test.isSentinel {
			t.Errorf("%s: %t, want %t", test.desc, got, test.isSentinel)
		}
	}
}

func TestCause(t *testing.T) {

	// set up a chain of errors
	const cause = etype.Sentinel("original error")
	err := xerrors.Errorf("first annotation: %w", cause)
	err = xerrors.Errorf("second annotation: %w", err)

	// expect Cause to find the very bottom error
	if !xerrors.Is(etype.Cause(err), cause) {
		t.Errorf("cause: %v, want %v", etype.Cause(err), cause)
	}
}
