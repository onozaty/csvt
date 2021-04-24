package cmd

import (
	"bufio"
	"bytes"
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

	table, err := loadCsvTable(r, "ID")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(table.columnNames(), []string{"ID", "Name", "Height", "Weight"}) {
		t.Fatal("failed test\n", table.columnNames())
	}

	if table.joinColumnName() != "ID" {
		t.Fatal("failed test\n", table.joinColumnName())
	}

	result := table.find("5")
	if !reflect.DeepEqual(
		result,
		map[string]string{
			"ID":     "5",
			"Name":   "Ichikawa",
			"Height": "152",
			"Weight": "50"}) {

		t.Fatal("failed test\n", result)
	}

	result = table.find("10")
	if result != nil {

		t.Fatal("failed test\n", result)
	}
}

func TestJoin(t *testing.T) {

	s1 := `ID,Name
1,Yamada
5,Ichikawa
2,"Hanako, Sato"
`
	r1 := csv.NewReader(strings.NewReader(s1))

	s2 := `ID,Height,Weight
1,171,50
2,160,60
5,152,50
`
	r2 := csv.NewReader(strings.NewReader(s2))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewWriter(w)

	err := join(r1, r2, "ID", out)

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	out.Flush()
	result := string(b.Bytes())

	expect := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
2,"Hanako, Sato",160,60
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestJoin_rightNone(t *testing.T) {

	s1 := `ID,Name
1,Yamada
5,Ichikawa
2,"Hanako, Sato"
`
	r1 := csv.NewReader(strings.NewReader(s1))

	s2 := `ID,Height,Weight
5,152,50
`
	r2 := csv.NewReader(strings.NewReader(s2))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewWriter(w)

	err := join(r1, r2, "ID", out)

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	out.Flush()
	result := string(b.Bytes())

	expect := `ID,Name,Height,Weight
1,Yamada,,
5,Ichikawa,152,50
2,"Hanako, Sato",,
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}
