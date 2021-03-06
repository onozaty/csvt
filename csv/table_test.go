package csv

import (
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// MemoryTable
func TestLoadCsvMemoryTable(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
2,"Hanako, Sato",160,60
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	table, err := LoadCsvMemoryTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer table.Close()

	if !reflect.DeepEqual(table.ColumnNames(), []string{"ID", "Name", "Height", "Weight"}) {
		t.Fatal("failed test\n", table.ColumnNames())
	}

	if table.KeyColumnName() != "ID" {
		t.Fatal("failed test\n", table.KeyColumnName())
	}

	result, err := table.Find("5")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(
		result,
		map[string]string{
			"ID":     "5",
			"Name":   "Ichikawa",
			"Height": "152",
			"Weight": "50"}) {

		t.Fatal("failed test\n", result)
	}

	result, err = table.Find("10")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if result != nil {
		t.Fatal("failed test\n", result)
	}
}

func TestLoadCsvMemoryTable_duplicateKey(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
1,"Dup",160,60
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvMemoryTable(r, "ID")
	if err == nil || err.Error() != "ID:1 is duplicated" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemoryTable_keyColumnNotFound(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvMemoryTable(r, "id")
	if err == nil || err.Error() != "id is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemoryTable_empty(t *testing.T) {

	s := ""
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvMemoryTable(r, "ID")
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemoryTable_changeLineOnly(t *testing.T) {

	s := "\n"
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvMemoryTable(r, "ID")
	if err == nil || err.Error() != "ID is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvMemoryTable_big(t *testing.T) {

	const maxId = 100000

	s := [maxId]string{}
	s[0] = "ID,Name,Age"
	for i := 1; i < maxId; i++ {
		s[i] = strconv.Itoa(i) + ",ABCDEFGHIJ,10"
	}

	r := NewCsvReader(strings.NewReader(strings.Join(s[:], "\n")), Format{})

	table, err := LoadCsvMemoryTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer table.Close()

	if !reflect.DeepEqual(table.ColumnNames(), []string{"ID", "Name", "Age"}) {
		t.Fatal("failed test\n", table.ColumnNames())
	}

	if table.KeyColumnName() != "ID" {
		t.Fatal("failed test\n", table.KeyColumnName())
	}

	for i := 1; i < maxId; i++ {
		id := strconv.Itoa(i)
		result, err := table.Find(id)
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if result["ID"] != id {
			t.Fatal("failed test\n", result["ID"])
		}
	}

	result, err := table.Find(strconv.Itoa(maxId))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if result != nil {
		t.Fatal("failed test\n", result)
	}
}

// FileTable
func TestLoadCsvFileTable(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
2,"Hanako, Sato",160,60
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	table, err := LoadCsvFileTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer table.Close()

	if !reflect.DeepEqual(table.ColumnNames(), []string{"ID", "Name", "Height", "Weight"}) {
		t.Fatal("failed test\n", table.ColumnNames())
	}

	if table.KeyColumnName() != "ID" {
		t.Fatal("failed test\n", table.KeyColumnName())
	}

	result, err := table.Find("5")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(
		result,
		map[string]string{
			"ID":     "5",
			"Name":   "Ichikawa",
			"Height": "152",
			"Weight": "50"}) {

		t.Fatal("failed test\n", result)
	}

	result, err = table.Find("10")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if result != nil {
		t.Fatal("failed test\n", result)
	}
}

func TestLoadCsvFileTable_duplicateKey(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
1,"Dup",160,60
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvFileTable(r, "ID")
	if err == nil || err.Error() != "ID:1 is duplicated" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileTable_keyColumnNotFound(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvFileTable(r, "id")
	if err == nil || err.Error() != "id is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileTable_empty(t *testing.T) {

	s := ""
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvFileTable(r, "ID")
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileTable_changeLineOnly(t *testing.T) {

	s := "\n"
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadCsvFileTable(r, "ID")
	if err == nil || err.Error() != "ID is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestLoadCsvFileTable_big(t *testing.T) {

	const maxId = 100000

	s := [maxId]string{}
	s[0] = "ID,Name,Age"
	for i := 1; i < maxId; i++ {
		s[i] = strconv.Itoa(i) + ",ABCDEFGHIJ,10"
	}

	r := NewCsvReader(strings.NewReader(strings.Join(s[:], "\n")), Format{})

	table, err := LoadCsvFileTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer table.Close()

	if !reflect.DeepEqual(table.ColumnNames(), []string{"ID", "Name", "Age"}) {
		t.Fatal("failed test\n", table.ColumnNames())
	}

	if table.KeyColumnName() != "ID" {
		t.Fatal("failed test\n", table.KeyColumnName())
	}

	for i := 1; i < maxId; i++ {
		id := strconv.Itoa(i)
		result, err := table.Find(id)
		if err != nil {
			t.Fatal("failed test\n", err)
		}

		if result["ID"] != id {
			t.Fatal("failed test\n", result["ID"])
		}
	}

	result, err := table.Find(strconv.Itoa(maxId))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if result != nil {
		t.Fatal("failed test\n", result)
	}
}
