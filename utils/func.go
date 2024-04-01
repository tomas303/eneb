package utils

type Result[V any, E any] struct {
	Value V
	Err   E
}
