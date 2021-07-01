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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-r", "a",
		"-t", "x",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col2",
		"-r", "a",
		"-t", "x",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-r", "a",
		"-t", "x",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-r", "^a$",
		"-t", "x",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-r", "^$",
		"-t", "xxx",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col2",
		"-r", ".*?([0-9]+).*",
		"-t", "#$1",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col2",
		"-r", `\s`,
		"-t", "",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

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
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col2",
		"-r", "a",
		"-t", "",
		"--delim", "\t",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

	expect := "id\tcol1\tcol2\r\n" +
		"1		bbb\r\n" +
		"2	bc	z\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestReplaceCmd_regex_invalid(t *testing.T) {

	s := ""

	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-r", `[a-`,
		"-t", "",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "regular expression specified in --regex is invalid: error parsing regexp: missing closing ]: `[a-`" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_fileNotFound(t *testing.T) {

	fi, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name() + "____", // 存在しないファイル名を指定
		"-o", fo.Name(),
		"-c", "col1",
		"-r", "a",
		"-t", "",
	})

	err = rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fi.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_columnNotFound(t *testing.T) {

	s := `id,col1,col2
1,a,b
2,c,
`

	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "colx", // 存在しないカラム
		"-r", "a",
		"-t", "",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "missing colx in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestReplaceCmd_empty(t *testing.T) {

	s := ""

	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"replace",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "colx",
		"-r", "a",
		"-t", "",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
