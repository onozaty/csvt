package cmd

import (
	"os"
	"testing"
)

func TestFilterCmd(t *testing.T) {

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

func TestFilterCmd_custom(t *testing.T) {

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
