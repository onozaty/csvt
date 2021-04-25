package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"strings"
	"testing"
)

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
