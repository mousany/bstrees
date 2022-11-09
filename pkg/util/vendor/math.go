package vendor

import "bstrees/pkg/trait/number"

func Min[T number.Number](args ...T) T {
	min := args[0]
	for _, arg := range args {
		if arg < min {
			min = arg
		}
	}
	return min
}

func Max[T number.Number](args ...T) T {
	max := args[0]
	for _, arg := range args {
		if arg > max {
			max = arg
		}
	}
	return max
}
