package utils

import "slices"

func Reverse[T any](collection []T) []T {
	slices.Reverse(collection)
	return collection
}

// Last returns the last element of a collection or error if empty.
func Last[T any](collection []T) (T, bool) {
	length := len(collection)

	if length == 0 {
		var t T
		return t, false
	}

	return collection[length-1], true
}

// LastOrEmpty returns the last element of a collection or zero value if empty.
func LastOrEmpty[T any](collection []T) T {
	i, _ := Last(collection)
	return i
}
