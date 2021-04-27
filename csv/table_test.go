package csv

import (
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

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
	if err == nil || err.Error() != "ID:1 is duplicated" {
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
	if err == nil || err.Error() != "id is not found" {
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

func TestLoadCsvTable_big(t *testing.T) {

	const maxId = 1000000

	s := [maxId]string{}
	s[0] = "ID,Name,Age"
	for i := 1; i < maxId; i++ {
		s[i] = strconv.Itoa(i) + ",ABCDEFGHIJ,10"
	}

	r, err := NewCsvReader(strings.NewReader(strings.Join(s[:], "\n")))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	table, err := LoadCsvTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(table.ColumnNames(), []string{"ID", "Name", "Age"}) {
		t.Fatal("failed test\n", table.ColumnNames())
	}

	if table.JoinColumnName() != "ID" {
		t.Fatal("failed test\n", table.JoinColumnName())
	}

	for i := 1; i < maxId; i++ {
		id := strconv.Itoa(i)
		result := table.Find(id)
		if result["ID"] != id {
			t.Fatal("failed test\n", result["ID"])
		}
	}

	result := table.Find(strconv.Itoa(maxId))
	if result != nil {
		t.Fatal("failed test\n", result)
	}
}
