package cmd

import (
	"os"
	"testing"
)

func TestFilterCmd(t *testing.T) {

	s := `ID,Name,CompanyID
,Yamada,
,"",
2,,""
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
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

	expect := "ID,Name,CompanyID\r\n" +
		",Yamada,\r\n" +
		"2,,\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_format(t *testing.T) {

	s := "ID;Name;CompanyID|1;Yamada;1|5;Ichikawa;|2;'Hanako; Sato';"
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"--delim", ";",
		"--quote", "'",
		"--sep", "|",
		"--allquote",
		"--bom",
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

	expect := "\uFEFF'ID';'Name';'CompanyID'|'1';'Yamada';'1'|"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_column(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",""
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
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

	expect := "ID,Name,CompanyID\r\n" +
		"1,Yamada,1\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_multiColumn(t *testing.T) {

	s := `col1,col2,col3
a,,
,b,
,,c
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col3",
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
		"a,,\r\n" +
		",,c\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equal(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "ID",
		"--equal", "1",
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

	expect := "ID,Name,CompanyID\r\n" +
		"1,Yamada,1\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equal_multiColumn(t *testing.T) {

	s := `col1,col2,col3
a,b,c
b,c,a
c,a,b
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col3",
		"--equal", "a",
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
		"a,b,c\r\n" +
		"b,c,a\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equal_allColumn(t *testing.T) {

	s := `col1,col2,col3
a,b,c
b,c,a
b,b,b
c,a,b
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--equal", "a",
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
		"a,b,c\r\n" +
		"b,c,a\r\n" +
		"c,a,b\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_regex(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, yamada",3
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "Name",
		"--regex", "[yY]amada",
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

	expect := "ID,Name,CompanyID\r\n" +
		"1,Yamada,1\r\n" +
		"2,\"Hanako, yamada\",3\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equal_regex_multiColumn(t *testing.T) {

	s := `col1,col2,col3
Ab,bc,
b,c,a
ba,a,b
abb,,ab
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "col1",
		"-c", "col2",
		"--regex", "(?i)^ab?$",
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
		"Ab,bc,\r\n" +
		"ba,a,b\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equal_regex_allColumn(t *testing.T) {

	s := `col1,col2,col3
ab,a,c
a,b,c
a,ba,bc
,,bb
a,,
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--regex", "b$",
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
		"ab,a,c\r\n" +
		"a,b,c\r\n" +
		",,bb\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_regex_invalid(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, yamada",3
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "Name",
		"--regex", "[a-z",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "regular expression specified in --regex is invalid: error parsing regexp: missing closing ]: `[a-z`" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_equal_regex(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, yamada",3
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "Name",
		"--equal", "Yamada",
		"--regex", "[a-z]",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --equal and --regex" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_fileNotFound(t *testing.T) {

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
		"filter",
		"-i", fi.Name() + "____", // 存在しないファイル名を指定
		"-o", fo.Name(),
		"-c", "CompanyID",
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

func TestFilterCmd_columnNotFound(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "Company", // 存在しないカラム
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_empty(t *testing.T) {

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
		"filter",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
