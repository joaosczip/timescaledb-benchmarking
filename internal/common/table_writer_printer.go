package common

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

type TableWriterPrinter struct {
	table *tablewriter.Table
}

func NewTableWriterPrinter() *TableWriterPrinter {
	return &TableWriterPrinter{
		table: tablewriter.NewWriter(os.Stdout),
	}
}

func (t *TableWriterPrinter) Print(data map[string]string) error {
	headers := []string{}
	rows := []string{}

	for k, v := range data {
		headers = append(headers, k)
		rows = append(rows, v)
	}

	t.table.SetHeader(headers)
	t.table.Append(rows)
	t.table.Render()

	return nil
}
