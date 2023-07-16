package bstrees

import (
	"errors"
)

var (
	ErrIndexIsOutOfRange       = errors.New("index is out of range")
	ErrPredecessorDoesNotExist = errors.New("predecessor does not exist")
	ErrSuccessorDoesNotExist   = errors.New("successor does not exist")
)
