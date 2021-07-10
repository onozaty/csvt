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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID\r\n" +
		",Yamada,\r\n" +
		"2,,\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_format(t *testing.T) {

	s := "ID;Name;CompanyID|1;Yamada;1|5;Ichikawa;|2;'Hanako; Sato';"
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "CompanyID",
		"--delim", ";",
		"--quote", "'",
		"--sep", "|",
		"--allquote",
		"--bom",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "CompanyID",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "ID",
		"--equal", "1",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col3",
		"--equal", "a",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"--equal", "a",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--regex", "[yY]amada",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID\r\n" +
		"1,Yamada,1\r\n" +
		"2,\"Hanako, yamada\",3\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_regex_multiColumn(t *testing.T) {

	s := `col1,col2,col3
Ab,bc,
b,c,a
ba,a,b
abb,,ab
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col2",
		"--regex", "(?i)^ab?$",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "col1,col2,col3\r\n" +
		"Ab,bc,\r\n" +
		"ba,a,b\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_regex_allColumn(t *testing.T) {

	s := `col1,col2,col3
ab,a,c
a,b,c
a,ba,bc
,,bb
a,,
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"--regex", "b$",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "col1,col2,col3\r\n" +
		"ab,a,c\r\n" +
		"a,b,c\r\n" +
		",,bb\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equalColumn(t *testing.T) {

	s := `col1,col2,col3
a,b,a
b,c,a
b,b,b
a,a,b
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"--column", "col1",
		"--equal-column", "col3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "col1,col2,col3\r\n" +
		"a,b,a\r\n" +
		"b,b,b\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_equalColumn_multiColumn(t *testing.T) {

	s := `col1,col2,col3
a,b,a
a,b,b
b,b,b
a,a,b
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"--column", "col1",
		"--column", "col2",
		"--equal-column", "col3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "col1,col2,col3\r\n" +
		"a,b,a\r\n" +
		"a,b,b\r\n" +
		"b,b,b\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestFilterCmd_not(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "ID",
		"--equal", "1",
		"--not",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID\r\n" +
		"5,Ichikawa,1\r\n" +
		"2,\"Hanako, Sato\",3\r\n"

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
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--regex", "[a-z",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "regular expression specified in --regex is invalid: error parsing regexp: missing closing ]: `[a-z`" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_equal_regex(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--equal", "Yamada",
		"--regex", "[a-z]",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --equal and --regex and --equal-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_equal_equalColumn(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--equal", "A",
		"--equal-column", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --equal and --regex and --equal-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_regex_equalColumn(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--regex", "A",
		"--equal-column", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --equal and --regex and --equal-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_equal_regex_equalColumn(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--equal", "A",
		"--regex", "A",
		"--equal-column", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --equal and --regex and --equal-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_fileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi + "____", // 存在しないファイル名を指定
		"-o", fo,
		"-c", "CompanyID",
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

func TestFilterCmd_columnNotFound(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Company", // 存在しないカラム
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_equalColumn_notFound(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "Name",
		"--equal-column", "Company", // 存在しないカラム
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "CompanyID",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestFilterCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"filter",
		"-i", fi,
		"-o", fo,
		"-c", "CompanyID",
		"--encoding", "xxxx",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "invalid encoding name: xxxx" {
		t.Fatal("failed test\n", err)
	}
}
