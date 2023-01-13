package errors

import (
	errorsPkg "errors"
)

var (
	ErrOutOfRange  = errorsPkg.New("index is out of range")
	ErrNoPrevValue = errorsPkg.New("previous value does not exist")
	ErrNoNextValue = errorsPkg.New("next value does not exist")

	ErrViolatedBSTProperty = errorsPkg.New("binary search tree property is violated")
	ErrViolatedRBProperty  = errorsPkg.New("red-black tree property is violated")
)
