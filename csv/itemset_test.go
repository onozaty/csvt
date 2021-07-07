package csv

import (
	"strings"
	"testing"
)

func TestNewItemSet(t *testing.T) {

	itemset := NewItemSet()
	if itemset.Count() != 0 {
		t.Fatal("failed test\n", itemset.Count())
	}
	if itemset.Contains("aa") {
		t.Fatal("failed test\n")
	}

	itemset.Add("aa")
	if itemset.Count() != 1 {
		t.Fatal("failed test\n", itemset.Count())
	}
	if !itemset.Contains("aa") {
		t.Fatal("failed test\n")
	}
	if itemset.Contains("a") {
		t.Fatal("failed test\n")
	}

	// 同じものを追加
	itemset.Add("aa")
	if itemset.Count() != 1 { // 数は増えない
		t.Fatal("failed test\n", itemset.Count())
	}
	if !itemset.Contains("aa") {
		t.Fatal("failed test\n")
	}
	if itemset.Contains("a") {
		t.Fatal("failed test\n")
	}

	itemset.Add("a")
	if itemset.Count() != 2 {
		t.Fatal("failed test\n", itemset.Count())
	}
	if !itemset.Contains("aa") {
		t.Fatal("failed test\n")
	}
	if !itemset.Contains("a") {
		t.Fatal("failed test\n")
	}
}

func TestLoadItemSet(t *testing.T) {

	s := `col1,col2
1,2
2,3
3,3
4,1
`

	r := NewCsvReader(strings.NewReader(s), Format{})

	itemset, err := LoadItemSet(r, "col2")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if itemset.Count() != 3 {
		t.Fatal("failed test\n", itemset.Count())
	}
	if !itemset.Contains("1") {
		t.Fatal("failed test\n")
	}
	if !itemset.Contains("2") {
		t.Fatal("failed test\n")
	}
	if !itemset.Contains("3") {
		t.Fatal("failed test\n")
	}
	if itemset.Contains("4") {
		t.Fatal("failed test\n")
	}
}

func TestLoadItemSet_columnNotFound(t *testing.T) {

	s := `col1,col2
1,2
`
	r := NewCsvReader(strings.NewReader(s), Format{})

	_, err := LoadItemSet(r, "col3")
	if err == nil || err.Error() != "col3 is not found" {
		t.Fatal("failed test\n", err)
	}
}
