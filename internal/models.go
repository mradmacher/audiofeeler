package internal

type FindResult[T any] struct {
	Record  T
	IsFound bool
}
