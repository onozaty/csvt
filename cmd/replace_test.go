package cmd

import (
	"os"
	"testing"
)

func TestReplaceCmd(t *testing.T) {

	s := `id,col1,col2
1,abc,abc
2,a,
3,aa  aa,A
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-r", "a",
		"-t", "x",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id,col1,col2\r\n" +
		"1,xbc,abc\r\n" +
		"2,x,\r\n" +
		"3,xx  xx,A\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_multiColumn(t *testing.T) {

	s := `id,col1,col2
1,abc,abc
2,a,
3,aa  aa,A
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col2",
		"-r", "a",
		"-t", "x",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id,col1,col2\r\n" +
		"1,xbc,xbc\r\n" +
		"2,x,\r\n" +
		"3,xx  xx,A\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_allColumn(t *testing.T) {

	s := `col1,col2,col3
z,abc,abc
a,a,
x,aa  aa,A
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-r", "a",
		"-t", "x",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "col1,col2,col3\r\n" +
		"z,xbc,xbc\r\n" +
		"x,x,\r\n" +
		"x,xx  xx,A\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_regex_full(t *testing.T) {

	s := `id,col1,col2
1,"a
bc",abc
2,a,
3,aa  aa,A
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-r", "^a$",
		"-t", "x",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id,col1,col2\r\n" +
		"1,\"a\nbc\",abc\r\n" +
		"2,x,\r\n" +
		"3,aa  aa,A\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_regex_empty(t *testing.T) {

	s := `id,col1,col2
1,abc,abc
2,,
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-r", "^$",
		"-t", "xxx",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id,col1,col2\r\n" +
		"1,abc,abc\r\n" +
		"2,xxx,\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_regex_capture(t *testing.T) {

	s := `id,col1,col2
1,aa123,xx2xx
2,9,
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col2",
		"-r", ".*?([0-9]+).*",
		"-t", "#$1",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id,col1,col2\r\n" +
		"1,#123,#2\r\n" +
		"2,#9,\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_regex_meta(t *testing.T) {

	s := `id,col1,col2
1,"a
 b",a   b
2,   ,
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col2",
		"-r", `\s`,
		"-t", "",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id,col1,col2\r\n" +
		"1,ab,ab\r\n" +
		"2,,\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_format(t *testing.T) {

	s := `id	col1	col2
1	aaa	bbb
2	abc	za
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col2",
		"-r", "a",
		"-t", "",
		"--delim", "\t",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "id\tcol1\tcol2\r\n" +
		"1		bbb\r\n" +
		"2	bc	z\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_regex_invalid(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-r", `[a-`,
		"-t", "",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "regular expression specified in --regex is invalid: error parsing regexp: missing closing ]: `[a-`" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_fileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi + "____", // 存在しないファイル名を指定
		"-o", fo,
		"-c", "col1",
		"-r", "a",
		"-t", "",
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fi+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_columnNotFound(t *testing.T) {

	s := `id,col1,col2
1,a,b
2,c,
`

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "colx", // 存在しないカラム
		"-r", "a",
		"-t", "",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing colx in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "colx",
		"-r", "a",
		"-t", "",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi,
		"-o", fo,
		"-c", "colx",
		"-r", "a",
		"-t", "",
		"--quote", "__",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag quote should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}
