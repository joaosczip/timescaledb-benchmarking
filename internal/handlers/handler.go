package handlers

type Handler[T any, R any] interface {
	Handle(input []T) *R
}
