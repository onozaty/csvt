package csv

import (
	"io"

	"github.com/onozaty/go-customcsv"
)

type Format struct {
	Delimiter       rune
	Quote           rune
	RecordSeparator string
	AllQuotes       bool
	WithBom         bool
}

type CsvReader interface {
	Read() (record []string, err error)
}
type CsvWriter interface {
	Write(record []string) error
	Flush() error
}

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func NewCsvReader(r io.Reader, f Format) CsvReader {

	// ReaderではBOMは自動的に除去
	cr := customcsv.NewReader(r)
	if f.Delimiter != 0 {
		cr.Delimiter = f.Delimiter
	}
	if f.Quote != 0 {
		cr.Quote = f.Quote
	}
	if f.RecordSeparator != "" {
		cr.SpecialRecordSeparator = f.RecordSeparator
	}

	return cr
}

func NewCsvWriter(w io.Writer, f Format) CsvWriter {

	if f.WithBom {
		w.Write(utf8bom)
	}

	cw := customcsv.NewWriter(w)
	if f.Delimiter != 0 {
		cw.Delimiter = f.Delimiter
	}
	if f.Quote != 0 {
		cw.Quote = f.Quote
	}
	if f.RecordSeparator != "" {
		cw.RecordSeparator = f.RecordSeparator
	}
	cw.AllQuotes = f.AllQuotes

	return cw
}
