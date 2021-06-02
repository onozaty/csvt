package csv

import (
	"bufio"
	"encoding/csv"
	"io"
	"reflect"
)

type CsvReader interface {
	Read() (record []string, err error)
}
type CsvWriter interface {
	Write(record []string) error
	Flush()
}

var utf8bom = []byte{0xEF, 0xBB, 0xBF}

func NewCsvReader(r io.Reader) (CsvReader, error) {

	br := bufio.NewReader(r)
	mark, err := br.Peek(len(utf8bom))

	if err != io.EOF && err != nil {
		return nil, err
	}

	if reflect.DeepEqual(mark, utf8bom) {
		// BOMがあれば読み飛ばす
		br.Discard(len(utf8bom))
	}

	return csv.NewReader(br), nil
}

func NewCsvWriter(w io.Writer) CsvWriter {

	cw := csv.NewWriter(w)
	// https://datatracker.ietf.org/doc/html/rfc4180#section-2
	cw.UseCRLF = true
	return cw
}
