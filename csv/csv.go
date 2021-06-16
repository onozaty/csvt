package csv

import (
	"io"

	"github.com/onozaty/go-customcsv"
)

type CsvReader interface {
	Read() (record []string, err error)
}
type CsvWriter interface {
	Write(record []string) error
	Flush() error
}

func NewCsvReader(r io.Reader) CsvReader {

	return customcsv.NewReader(r)
}

func NewCsvWriter(w io.Writer) CsvWriter {

	return customcsv.NewWriter(w)
}
