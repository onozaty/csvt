package csv

import (
	"encoding/csv"
	"reflect"
	"strings"
	"testing"
)

func TestLoadCsvTable(t *testing.T) {

	s := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
2,"Hanako, Sato",160,60
`
	r := csv.NewReader(strings.NewReader(s))

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
