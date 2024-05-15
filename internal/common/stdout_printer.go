package common

type StdoutPrinter interface {
	Print(data map[string]string) error
}
