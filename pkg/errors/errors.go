package errors

import (
	defaultErrors "errors"
)

var (
	ErrOutOfRange  = defaultErrors.New("out of range")
	ErrNoPrevValue = defaultErrors.New("no previous value")
	ErrNoNextValue = defaultErrors.New("no next value")
)
