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
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "3\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestCountCmd_format(t *testing.T) {

	s := `ID	Name	CompanyID
1	Yamada	1
5	Ichikawa	
2	"Hanako	Sato"	3
`
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
		"--delim", `\t`,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
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
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
		"-c", "CompanyID",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
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
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
		"--header",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if result != "4\n" {
		t.Fatal("failed test\n", result)
	}
}

func TestCountCmd_fileNotFound(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
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

func TestCountCmd_columnNotFound(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,
2,"Hanako, Sato",3
`
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
		"-c", "Company", // 存在しないカラム
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestCountCmd_empty(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestCountCmd_args(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
		"aaaa", // フラグ以外を指定
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "only flags can be specified" {
		t.Fatal("failed test\n", err)
	}
}

func TestCountCmd_invalidFormat(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"count",
		"-i", f,
		"--delim", "aa",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}
