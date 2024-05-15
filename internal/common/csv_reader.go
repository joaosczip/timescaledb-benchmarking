package common

type CsvReader[T any] interface {
	Read(path string) ([]T, error)
}
