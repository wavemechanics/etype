package etype

// A Sentinel is a simple error string that can be a const.
// See https://dave.cheney.net/2019/06/10/constant-time
type Sentinel string

func (e Sentinel) Error() string {
	return string(e)
}

// Is returns true if err is a Sentinel and is the same kind of sentinel as the receiver.
// (ie, the string values are the same)
func (e Sentinel) Is(err error) bool {
	sentinel, ok := err.(Sentinel)
	if !ok {
		return false
	}
	return sentinel == e
}

// Cause returns the lowest error in the chain.
func Cause(err error) error {
	type wrapper interface {
		Unwrap() error
	}
	for err != nil {
		cause, ok := err.(wrapper)
		if !ok {
			break
		}
		err = cause.Unwrap()
	}
	return err
}
