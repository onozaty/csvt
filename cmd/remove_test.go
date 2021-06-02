package cmd

import (
	"os"
	"testing"
)

func TestRemoveCmd(t *testing.T) {

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
		"remove",
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

	expect := "ID,Name\r\n" +
		"1,Yamada\r\n" +
		"5,Ichikawa\r\n" +
		"2,\"Hanako, Sato\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRemoveCmd_columns(t *testing.T) {

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
		"remove",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"-c", "ID",
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

	expect := "Name\r\n" +
		"Yamada\r\n" +
		"Ichikawa\r\n" +
		"\"Hanako, Sato\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRemoveCmd_fileNotFound(t *testing.T) {

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
		"remove",
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

func TestRemoveCmd_columnNotFound(t *testing.T) {

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
		"remove",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "Company", // 存在しないカラム
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestRemoveCmd_empty(t *testing.T) {

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
		"remove",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
