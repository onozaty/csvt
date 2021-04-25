package csv

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestNewCsvReader_withBOM(t *testing.T) {

	// 先頭にBOM
	s := "\uFEFFID,Name\n1,Yamada"

	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	header, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// BOMが先頭についていないこと
	if !reflect.DeepEqual(header, []string{"ID", "Name"}) {
		t.Fatal("failed test\n", header)
	}

	firstRow, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	if !reflect.DeepEqual(firstRow, []string{"1", "Yamada"}) {
		t.Fatal("failed test\n", header)
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewCsvReader_withoutBOM(t *testing.T) {

	s := "ID,Name\n1,Yamada"

	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	header, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// BOMが先頭についていないこと
	if !reflect.DeepEqual(header, []string{"ID", "Name"}) {
		t.Fatal("failed test\n", header)
	}

	firstRow, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	if !reflect.DeepEqual(firstRow, []string{"1", "Yamada"}) {
		t.Fatal("failed test\n", header)
	}

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewCsvWriter(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewCsvWriter(w)

	cw.Write([]string{"1", "2"})
	cw.Write([]string{"あ", "a"})
	cw.Write([]string{",", ""})

	cw.Flush()
	result := string(b.Bytes())

	expect := `1,2
あ,a
",",
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestLoadCsvTable(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
2,"Hanako, Sato",160,60
`
	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	table, err := LoadCsvTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(table.ColumnNames(), []string{"ID", "Name", "Height", "Weight"}) {
		t.Fatal("failed test\n", table.ColumnNames())
	}

	if table.JoinColumnName() != "ID" {
		t.Fatal("failed test\n", table.JoinColumnName())
	}

	result := table.Find("5")
	if !reflect.DeepEqual(
		result,
		map[string]string{
			"ID":     "5",
			"Name":   "Ichikawa",
			"Height": "152",
			"Weight": "50"}) {

		t.Fatal("failed test\n", result)
	}

	result = table.Find("10")
	if result != nil {

		t.Fatal("failed test\n", result)
	}
}

func TestLoadCsvTable_duplicateKey(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
1,"Dup",160,60
`
	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = LoadCsvTable(r, "ID")
	if err == nil || err.Error() != "ID:1 is duplicated." {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvTable_joinColumnNotFound(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
`
	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = LoadCsvTable(r, "id")
	if err == nil || err.Error() != "id is not found." {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvTable_empty(t *testing.T) {

	s := ""
	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = LoadCsvTable(r, "ID")
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvTable_changeLineOnly(t *testing.T) {

	s := "\n"
	r, err := NewCsvReader(strings.NewReader(s))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	_, err = LoadCsvTable(r, "ID")
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}
