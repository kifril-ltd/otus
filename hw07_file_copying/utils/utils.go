package utils

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		return Zero[T]()
	}

	min := args[0]
	for _, arg := range args[1:] {
		if arg < min {
			min = arg
		}
	}
	return min
}

func Zero[T any]() T {
	var zero T
	return zero
}
