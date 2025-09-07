package slices

import (
	"slices"
)

func Contains[S ~[]E, E comparable](s S, v E) bool {
	return slices.Contains(s, v)
}
