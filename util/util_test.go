package util

import (
	"reflect"
	"testing"
)

func TestIndexOf(t *testing.T) {

	strings := []string{"a", "b", "c"}

	index := IndexOf(strings, "a")
	if index != 0 {
		t.Fatal("failed test\n", index)
	}

	index = IndexOf(strings, "b")
	if index != 1 {
		t.Fatal("failed test\n", index)
	}

	index = IndexOf(strings, "c")
	if index != 2 {
		t.Fatal("failed test\n", index)
	}

	index = IndexOf(strings, "d")
	if index != -1 {
		t.Fatal("failed test\n", index)
	}
}

func TestIndexOf_multi(t *testing.T) {

	strings := []string{"a", "a", "a"}

	// 最初に見つかった位置になること
	index := IndexOf(strings, "a")
	if index != 0 {
		t.Fatal("failed test\n", index)
	}
}

func TestRemove(t *testing.T) {

	removed := Remove([]string{"a", "b", "c"}, "a")
	if !reflect.DeepEqual(removed, []string{"b", "c"}) {
		t.Fatal("failed test\n", removed)
	}

	removed = Remove([]string{"a", "b", "c"}, "b")
	if !reflect.DeepEqual(removed, []string{"a", "c"}) {
		t.Fatal("failed test\n", removed)
	}

	removed = Remove([]string{"a", "b", "c"}, "c")
	if !reflect.DeepEqual(removed, []string{"a", "b"}) {
		t.Fatal("failed test\n", removed)
	}
}

func TestRemove_multi(t *testing.T) {

	removed := Remove([]string{"a", "b", "c", "b"}, "b")
	if !reflect.DeepEqual(removed, []string{"a", "c"}) {
		t.Fatal("failed test\n", removed)
	}
}

func TestContains(t *testing.T) {

	array := []int{10, 11, 20}

	if !Contains(array, 10) {
		t.Fatal("failed test\n")
	}

	if !Contains(array, 11) {
		t.Fatal("failed test\n")
	}

	if !Contains(array, 20) {
		t.Fatal("failed test\n")
	}

	if Contains(array, 2) {
		t.Fatal("failed test\n")
	}
}
