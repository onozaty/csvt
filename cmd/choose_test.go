package cmd

import (
	"os"
	"testing"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func TestChooseCmd(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := "CompanyID\r\n" +
		"1\r\n" +
		"1\r\n" +
		"3\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestChooseCmd_format(t *testing.T) {

	s := "ID;Name;CompanyID|1;Yamada;1|5;Ichikawa;1|2;'Hanako; Sato';3"
	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"--delim", ";",
		"--quote", "'",
		"--sep", "|",
		"--allquote",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

	expect := "'CompanyID'|'1'|'1'|'3'|"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestChooseCmd_columns(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"-c", "ID",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := "ID,CompanyID\r\n" +
		"1,1\r\n" +
		"5,1\r\n" +
		"2,3\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestChooseCmd_fileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name() + "____", // 存在しないファイル名を指定
		"-o", fo.Name(),
		"-c", "CompanyID",
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fi.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestChooseCmd_columnNotFound(t *testing.T) {

	s := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "Company", // 存在しないカラム
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing Company in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestChooseCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestChooseCmd_encoding(t *testing.T) {

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", "../testdata/users-sjis.csv",
		"-o", fo.Name(),
		"-c", "名前",
		"--encoding", "shift_jis",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	b, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	// Shift_JIS->UTF-8に変換して期待値と比較
	ub, _, _ := transform.Bytes(japanese.ShiftJIS.NewDecoder(), b)
	result := string(ub)

	expect := "名前\r\n" +
		"\"Taro, Yamada\"\r\n" +
		"山田 花子\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestChooseCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"choose",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"--quote", "xxx",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag quote should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}
