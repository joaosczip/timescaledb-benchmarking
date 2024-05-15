package common

type StdoutWriter interface {
	Write(data map[string]string) error
}
