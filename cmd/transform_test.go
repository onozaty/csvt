package cmd

import (
	"os"
	"reflect"
	"testing"
)

func TestTransformCmd(t *testing.T) {

	s := `ID,Name
1,Taro; Yamada
2,"Hanako, Sato"
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--out-delim", ";",
		"--out-quote", "'",
		"--out-sep", "|",
		"--out-allquote",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "'ID';'Name'|'1';'Taro; Yamada'|'2';'Hanako, Sato'|"
	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_format(t *testing.T) {

	s := "ID/Name%1/Taro; Yamada%2/$Hanako, Sato$%"

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--delim", "/",
		"--quote", "$",
		"--sep", "%",
		"--out-delim", ";",
		"--out-quote", "'",
		"--out-sep", "|",
		"--out-allquote",
		"--out-bom",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "\uFEFF'ID';'Name'|'1';'Taro; Yamada'|'2';'Hanako, Sato'|"
	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_fileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi + "____", // 存在しないファイル名を指定
		"-o", fo,
		"--out-allquote",
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

func TestTransformCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--out-allquote",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	if result != "" {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_delim(t *testing.T) {

	s := `ID	Name
1	Taro
2	"Hanako	Sato"
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,Hanako	Sato\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_delim_multichar(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--delim", `;;`,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_delim_multibyte(t *testing.T) {

	s := `ID　Name
1　Taro
2　"Hanako　Sato"
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--delim", `\u3000`, // マルチバイトとなる全角スペース
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,Hanako　Sato\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_delim_parseError(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--delim", `\t"`,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != `Could not parse value \t" of flag delim: invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_quote(t *testing.T) {

	s := `'ID','Name'
'1',Taro
2,'Hanako, Sato'
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--quote", "'",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,\"Hanako, Sato\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_quote_multichar(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--quote", "''",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag quote should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_sep(t *testing.T) {

	s := `ID,Name|1,Taro|2,"Hanako, Sato"`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--sep", "|",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,\"Hanako, Sato\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_sep_parseError(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--sep", `\r"`,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != `Could not parse value \r" of flag sep: invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_encoding(t *testing.T) {

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", "../testdata/users-sjis.csv",
		"-o", fo,
		"--encoding", "sjis",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,名前,年齢\r\n" +
		"1,\"Taro, Yamada\",20\r\n" +
		"2,山田 花子,21\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_out_encoding(t *testing.T) {

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", "../testdata/users-utf8.csv",
		"-o", fo,
		"--out-encoding", "sjis",
		"--out-bom", // UTF-8ではないのでBOM指定しても付かない
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result, err := os.ReadFile(fo)
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	expect, err := os.ReadFile("../testdata/users-sjis.csv")
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_encoding_invalid(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi,
		"-o", fo,
		"--out-encoding", "utf",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "invalid encoding name: utf" {
		t.Fatal("failed test\n", err)
	}
}
