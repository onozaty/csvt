package csv

import (
	"fmt"
	"io"
	"strings"

	"github.com/onozaty/go-customcsv"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type Format struct {
	Delimiter       rune
	Quote           rune
	RecordSeparator string
	AllQuotes       bool
	WithBom         bool
	Encoding        encoding.Encoding
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

	if f.Encoding != nil {
		r = transform.NewReader(r, f.Encoding.NewDecoder())
	}

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

	if f.Encoding != nil {
		w = transform.NewWriter(w, f.Encoding.NewEncoder())
	}

	if f.WithBom && f.Encoding == nil { // UTF-8の場合のみ
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

func Encoding(name string) (encoding.Encoding, error) {

	normalizeName := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(name), "_", ""), "-", "")

	switch normalizeName {
	case "utf8":
		// UTF-8の場合は変換不要
		return nil, nil
	case "shiftjis", "sjis":
		return japanese.ShiftJIS, nil
	case "eucjp":
		return japanese.EUCJP, nil
	default:
		return nil, fmt.Errorf("invalid encoding name: %s", name)
	}

}
