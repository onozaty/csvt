package csv

import (
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// Memory
func TestLoadCsvMemorySortedRows(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"2", "b"},
		[]string{"5", "d"},
		[]string{"1", "c"},
		[]string{"3", "a"},
		[]string{"4", "e"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	rows, err := LoadCsvMemorySortedRows(r, []string{"col1"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 5 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"1", "c"},
		[]string{"2", "b"},
		[]string{"3", "a"},
		[]string{"4", "e"},
		[]string{"5", "d"},
	)
}

func TestLoadCsvMemorySortedRows_multiColumn(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "c"},
		[]string{"2", "a"},
		[]string{"1", "a"},
		[]string{"2", "b"},
		[]string{"1", "b"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	rows, err := LoadCsvMemorySortedRows(r, []string{"col1", "col2"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 5 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"1", "a"},
		[]string{"1", "b"},
		[]string{"1", "c"},
		[]string{"2", "a"},
		[]string{"2", "b"},
	)
}

func TestLoadCsvMemorySortedRows_num(t *testing.T) {

	s := joinRows(
		[]string{"col1"},
		[]string{"10"},
		[]string{"2"},
		[]string{"9"},
		[]string{"123"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	rows, err := LoadCsvMemorySortedRows(r, []string{"col1"}, CompareNumber)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 4 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"2"},
		[]string{"9"},
		[]string{"10"},
		[]string{"123"},
	)
}

func TestLoadCsvMemorySortedRows_same(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "3"},
		[]string{"2", "1"},
		[]string{"1", "1"},
		[]string{"2", "2"},
		[]string{"1", "2"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	// col1だけ指定して同じ値がどうなるか確認
	rows, err := LoadCsvMemorySortedRows(r, []string{"col1"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 5 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"1", "3"},
		[]string{"1", "1"},
		[]string{"1", "2"},
		[]string{"2", "1"},
		[]string{"2", "2"},
	)
}

func TestLoadCsvMemorySortedRows_empty(t *testing.T) {

	r := NewCsvReader(strings.NewReader(""), Format{})

	_, err := LoadCsvMemorySortedRows(r, []string{"col1"}, CompareString)

	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemorySortedRows_columnNotFound(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "3"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvMemorySortedRows(r, []string{"col1", "col3"}, CompareString)

	if err == nil || err.Error() != "col3 is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemorySortedRows_invalidNumber(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "1"},
		[]string{"a", "2"}, //  数字じゃない
		[]string{"3", "3"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvMemorySortedRows(r, []string{"col1"}, CompareNumber)

	if err == nil || err.Error() != `strconv.Atoi: parsing "a": invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemorySortedRows_big(t *testing.T) {

	const maxId = 100000

	s := [maxId + 1]string{}
	s[0] = "col1,col2"
	for i := 1; i <= maxId; i++ {
		s[i] = strconv.Itoa(i) + "," + strconv.Itoa(maxId-i)
	}

	r := NewCsvReader(strings.NewReader(strings.Join(s[:], "\n")), Format{})

	rows, err := LoadCsvMemorySortedRows(r, []string{"col2"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != maxId {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	// 先頭と末尾を確認
	{
		row, err := rows.Row(0) // 先頭
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(row, []string{strconv.Itoa(maxId), "0"}) {
			t.Fatal("failed test\n", row)
		}
	}
	{
		row, err := rows.Row(maxId - 1) // 末尾
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(row, []string{"1", strconv.Itoa(maxId - 1)}) {
			t.Fatal("failed test\n", row)
		}
	}
}

// File
func TestLoadCsvFileSortedRows(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"2", "b"},
		[]string{"5", "d"},
		[]string{"1", "c"},
		[]string{"3", "a"},
		[]string{"4", "e"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	rows, err := LoadCsvFileSortedRows(r, []string{"col1"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 5 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"1", "c"},
		[]string{"2", "b"},
		[]string{"3", "a"},
		[]string{"4", "e"},
		[]string{"5", "d"},
	)
}

func TestLoadCsvFileSortedRows_multiColumn(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "c"},
		[]string{"2", "a"},
		[]string{"1", "a"},
		[]string{"2", "b"},
		[]string{"1", "b"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	rows, err := LoadCsvFileSortedRows(r, []string{"col1", "col2"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 5 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"1", "a"},
		[]string{"1", "b"},
		[]string{"1", "c"},
		[]string{"2", "a"},
		[]string{"2", "b"},
	)
}

func TestLoadCsvFileSortedRows_num(t *testing.T) {

	s := joinRows(
		[]string{"col1"},
		[]string{"10"},
		[]string{"2"},
		[]string{"9"},
		[]string{"123"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	rows, err := LoadCsvFileSortedRows(r, []string{"col1"}, CompareNumber)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 4 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"2"},
		[]string{"9"},
		[]string{"10"},
		[]string{"123"},
	)
}

func TestLoadCsvFileSortedRows_same(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "3"},
		[]string{"2", "1"},
		[]string{"1", "1"},
		[]string{"2", "2"},
		[]string{"1", "2"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	// col1だけ指定して同じ値がどうなるか確認
	rows, err := LoadCsvFileSortedRows(r, []string{"col1"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != 5 {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	assertRows(t, rows,
		[]string{"1", "3"},
		[]string{"1", "1"},
		[]string{"1", "2"},
		[]string{"2", "1"},
		[]string{"2", "2"},
	)
}

func TestLoadCsvFileSortedRows_empty(t *testing.T) {

	r := NewCsvReader(strings.NewReader(""), Format{})

	_, err := LoadCsvFileSortedRows(r, []string{"col1"}, CompareString)

	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileSortedRows_columnNotFound(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "3"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvFileSortedRows(r, []string{"col1", "col3"}, CompareString)

	if err == nil || err.Error() != "col3 is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileSortedRows_invalidNumber(t *testing.T) {

	s := joinRows(
		[]string{"col1", "col2"},
		[]string{"1", "1"},
		[]string{"a", "2"}, //  数字じゃない
		[]string{"3", "3"},
	)

	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvFileSortedRows(r, []string{"col1"}, CompareNumber)

	if err == nil || err.Error() != `strconv.Atoi: parsing "a": invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileSortedRows_big(t *testing.T) {

	const maxId = 100000

	s := [maxId + 1]string{}
	s[0] = "col1,col2"
	for i := 1; i <= maxId; i++ {
		s[i] = strconv.Itoa(i) + "," + strconv.Itoa(maxId-i)
	}

	r := NewCsvReader(strings.NewReader(strings.Join(s[:], "\n")), Format{})

	rows, err := LoadCsvFileSortedRows(r, []string{"col2"}, CompareString)

	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer rows.Close()

	if rows.Count() != maxId {
		t.Fatal("failed test\n", rows.Count())
	}

	if !reflect.DeepEqual(rows.ColumnNames(), []string{"col1", "col2"}) {
		t.Fatal("failed test\n", rows.ColumnNames())
	}

	// 先頭と末尾を確認
	{
		row, err := rows.Row(0) // 先頭
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(row, []string{strconv.Itoa(maxId), "0"}) {
			t.Fatal("failed test\n", row)
		}
	}
	{
		row, err := rows.Row(maxId - 1) // 末尾
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(row, []string{"1", strconv.Itoa(maxId - 1)}) {
			t.Fatal("failed test\n", row)
		}
	}
}

func assertRows(t *testing.T, rows CsvSortedRows, expecteds ...[]string) {

	for i, expected := range expecteds {

		row, err := rows.Row(i)
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if !reflect.DeepEqual(row, expected) {
			t.Fatal("failed test\n", i, row)
		}
	}
}

func joinRows(rows ...[]string) string {

	csv := ""

	for _, row := range rows {
		csv += strings.Join(row, ",") + "\r\n"
	}

	return csv
}
