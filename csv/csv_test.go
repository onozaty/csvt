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

	r := NewCsvReader(strings.NewReader(s), Format{})

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

	r := NewCsvReader(strings.NewReader(s), Format{})

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

func TestNewCsvReader_LF_CRLF(t *testing.T) {

	s := "ID,Name\n1,Yamada\r\n2,Ichikawa"

	r := NewCsvReader(strings.NewReader(s), Format{})

	header, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

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

	secondRow, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	if !reflect.DeepEqual(secondRow, []string{"2", "Ichikawa"}) {
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
	cw := NewCsvWriter(w, Format{})

	cw.Write([]string{"1", "2"})
	cw.Write([]string{"あ", "a"})
	cw.Write([]string{",", ""})

	cw.Flush()
	result := b.String()

	expect := "1,2\r\n" +
		"あ,a\r\n" +
		"\",\",\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}
