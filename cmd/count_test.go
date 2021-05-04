package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestCountCmd(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",3
`
	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f.Name(),
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "3\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestCountCmd_column(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",3
`
	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f.Name(),
		"-c", "CompanyID",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "2\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestCountCmd_header(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",3
`
	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f.Name(),
		"--header",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "4\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestCountCmd_fileNotFound(t *testing.T) {

	f, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f.Name() + "____", // 存在しないファイル名を指定
	})

	err = rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != f.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestCountCmd_columnNotFound(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",3
`
	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f.Name(),
		"-c", "Company", // 存在しないカラム
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestCountCmd_empty(t *testing.T) {

	s := ""

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f.Name(),
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
