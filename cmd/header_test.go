package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestHeaderCmd(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",3
`
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\nName\nCompanyID\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_format(t *testing.T) {

	s := "ID;Name;CompanyID|1;Yamada;1|5;Ichikawa;1|2;'Hanako; Sato';3"
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f,
		"--delim", ";",
		"--quote", "'",
		"--sep", "|",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\nName\nCompanyID\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_fileNotFound(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f + "____", // 存在しないファイル名を指定
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != f+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestHeaderCmd_empty(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestHeaderCmd_encoding_default(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", "../testdata/users-utf8.csv",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\n名前\n年齢\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_encoding_utf8(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", "../testdata/users-utf8.csv",
		"--encoding", "utf-8",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\n名前\n年齢\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_encoding_sjis(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", "../testdata/users-sjis.csv",
		"--encoding", "sjis",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\n名前\n年齢\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_encoding_shift_jis(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", "../testdata/users-sjis.csv",
		"--encoding", "Shift_JIS",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\n名前\n年齢\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_encoding_eucjp(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", "../testdata/users-eucjp.csv",
		"--encoding", "eucjp",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\n名前\n年齢\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_encoding_euc_jp(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", "../testdata/users-eucjp.csv",
		"--encoding", "EUC-JP",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\n名前\n年齢\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_invalidFormat(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f,
		"--quote", "zz",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag quote should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}
