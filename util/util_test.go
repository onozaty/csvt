package util

import (
	"reflect"
	"testing"
)

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
