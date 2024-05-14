package common

import "errors"

var ErrInvalidHeader = errors.New("invalid csv header")

type CsvReader[T any] interface {
	Read(path string) (T, error)
}
