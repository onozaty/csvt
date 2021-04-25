package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"reflect"

	"github.com/onozaty/csvt/util"
)

// CSV Reader / Writer
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

	return csv.NewWriter(w)
}

// CSV Table

type CsvTable interface {
	Find(key string) map[string]string
	JoinColumnName() string
	ColumnNames() []string
}

type MemoryTable struct {
	joinColumnName string
	columnNames    []string
	rows           map[string][]string
}

func (t *MemoryTable) Find(key string) map[string]string {

	row := t.rows[key]

	if row == nil {
		return nil
	}

	rowMap := make(map[string]string)
	for i := 0; i < len(t.columnNames); i++ {
		rowMap[t.columnNames[i]] = row[i]
	}

	return rowMap
}

func (t *MemoryTable) JoinColumnName() string {

	return t.joinColumnName
}

func (t *MemoryTable) ColumnNames() []string {

	return t.columnNames
}

func LoadCsvTable(reader CsvReader, joinColumnName string) (CsvTable, error) {

	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	primaryColumnIndex := util.IndexOf(headers, joinColumnName)
	if primaryColumnIndex == -1 {
		return nil, fmt.Errorf("%s is not found.", joinColumnName)
	}

	rows := make(map[string][]string)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// 格納前に既にあるか確認
		// -> 重複して存在した場合はエラーに
		_, has := rows[row[primaryColumnIndex]]
		if has {
			return nil, fmt.Errorf("%s:%s is duplicated.", joinColumnName, row[primaryColumnIndex])
		}

		rows[row[primaryColumnIndex]] = row
	}

	return &MemoryTable{
		joinColumnName: joinColumnName,
		columnNames:    headers,
		rows:           rows,
	}, nil
}
