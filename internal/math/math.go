package math

import "bstrees/internal/node"

func Min[T node.Number](args ...T) T {
	min := args[0]
	for _, arg := range args {
		if arg < min {
			min = arg
		}
	}
	return min
}

func Max[T node.Number](args ...T) T {
	max := args[0]
	for _, arg := range args {
		if arg > max {
			max = arg
		}
	}
	return max
}

func Abs[T node.Number](value T) T {
	if value < 0 {
		return -value
	}
	return value
}
