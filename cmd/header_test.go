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
	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f.Name(),
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "ID\nName\nCompanyID\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestHeaderCmd_fileNotFound(t *testing.T) {

	f, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
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

func TestHeaderCmd_empty(t *testing.T) {

	s := ""

	f, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"header",
		"-i", f.Name(),
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
