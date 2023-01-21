package csv

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/text/encoding/japanese"
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

func TestNewCsvReader_format(t *testing.T) {

	s := "ID;Name|1;'Yamada'"

	r := NewCsvReader(strings.NewReader(s), Format{
		Delimiter:       ';',
		Quote:           '\'',
		RecordSeparator: "|",
	})

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

	_, err = r.Read()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestNewCsvReader_encoding(t *testing.T) {

	f, err := os.Open("../testdata/users-sjis.csv")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	encoding, _ := Encoding("sjis")
	r := NewCsvReader(f, Format{
		Encoding: encoding,
	})

	header, err := r.Read()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(header, []string{"ID", "名前", "年齢"}) {
		t.Fatal("failed test\n", header)
	}

	{
		row, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}
		if !reflect.DeepEqual(row, []string{"1", "Taro, Yamada", "20"}) {
			t.Fatal("failed test\n", header)
		}
	}

	{
		row, err := r.Read()
		if err != nil {
			t.Fatal("failed test\n", err)
		}
		if !reflect.DeepEqual(row, []string{"2", "山田 花子", "21"}) {
			t.Fatal("failed test\n", header)
		}
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

func TestNewCsvWriter_withBom(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewCsvWriter(w, Format{
		WithBom: true,
	})

	cw.Write([]string{"1", "2"})

	cw.Flush()
	result := b.String()

	expect := "\uFEFF1,2\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewCsvWriter_format(t *testing.T) {

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewCsvWriter(w, Format{
		Delimiter:       ';',
		Quote:           '\'',
		RecordSeparator: "|",
		AllQuotes:       true,
	})

	cw.Write([]string{"1", "2"})
	cw.Write([]string{"あ", "a"})
	cw.Write([]string{";", ""})

	cw.Flush()
	result := b.String()

	expect := "'1';'2'|'あ';'a'|';';''|"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestNewCsvWriter_encoding(t *testing.T) {

	encoding, _ := Encoding("sjis")

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	cw := NewCsvWriter(w, Format{
		Encoding: encoding,
		WithBom:  true, // UTF-8以外なのでBOM指定してもBOM付かない
	})

	cw.Write([]string{"あ", "a"})

	cw.Flush()
	result := b.Bytes()

	expect, _ := japanese.ShiftJIS.NewEncoder().Bytes([]byte("あ,a\r\n"))

	if reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestEncoding_sjis(t *testing.T) {

	encoding1, err := Encoding("sjis")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	encoding2, err := Encoding("Shift_JIS")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if encoding1 != encoding2 {
		t.Fatal("failed test\n")
	}
}

func TestEncoding_eucjp(t *testing.T) {

	encoding1, err := Encoding("euc-jp")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	encoding2, err := Encoding("EUC-JP")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if encoding1 != encoding2 {
		t.Fatal("failed test\n")
	}
}

func TestEncoding_utf8(t *testing.T) {

	encoding1, err := Encoding("utf8")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	encoding2, err := Encoding("UTF-8")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if encoding1 != nil || encoding2 != nil {
		t.Fatal("failed test\n")
	}
}

func TestEncoding_invalid(t *testing.T) {

	_, err := Encoding("xxxx")
	if err == nil || err.Error() != "invalid encoding name: xxxx" {
		t.Fatal("failed test\n", err)
	}
}
