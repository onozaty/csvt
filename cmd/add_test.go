package cmd

import (
	"os"
	"testing"
)

func TestAddCmd(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,a,",
		"2,b,",
		"3,c,",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_value(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--value", "x",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,a,x",
		"2,b,x",
		"3,c,x",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_value_useBackslash(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--value", `\u0040`, // バックスラッシュを含む値
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,a,@",
		"2,b,@",
		"3,c,@",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_copyColumn(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--copy-column", "col1",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,a,1",
		"2,b,2",
		"3,c,3",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_template(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", "{{.col1}}-{{.col2}}",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,a,1-a",
		"2,b,2-b",
		"3,c,3-c",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_template_useBackslash(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", `{{.col1}}\n{{.col2}}`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,a,\"1\na\"",
		"2,b,\"2\nb\"",
		"3,c,\"3\nc\"",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_format(t *testing.T) {

	s := joinRows(
		"col1\tcol2",
		"1\ta",
		"2\tb",
		"3\tc",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", "{{ if eq .col1 \"2\" }}{{ .col1 }}{{ else }}{{ .col2 }}{{ end }}",
		"--delim", "\t",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1\tcol2\tcol3",
		"1\ta\ta",
		"2\tb\t2",
		"3\tc\tc",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestAddCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--encoding", "aa",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "invalid encoding name: aa" {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_templateParseError(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", "{{.col1}", // 閉じ括弧が少ない
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != `--template is invalid: template: template:1: bad character U+007D '}'` {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_valueUnquoteError(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--value", `\n"`, // \が含まれる状態で、エスケープされていないバックスラッシュ
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != `Could not parse value \n" of flag value: invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_templateUnquoteError(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", `\n"`, // \が含まれる状態で、エスケープされていないバックスラッシュ
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != `Could not parse value \n" of flag template: invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_value_template(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--value", "a",
		"--template", "{{.col1}}",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --value and --template and --copy-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_value_copyColumn(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--value", "a",
		"--copy-column", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --value and --template and --copy-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_template_copyColumn(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", "{{.col1}}",
		"--copy-column", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "not allowed to specify both --value and --template and --copy-column" {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi + "____", // 存在しないファイル
		"-o", fo,
		"-c", "col3",
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

func TestAddCmd_copyColumnNotFound(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--copy-column", "col3",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing col3 in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_inputFileEmpty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--copy-column", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestAddCmd_templateExecuteError(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"add",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
		"--template", "{{if eq 1 \"1\"}}xx{{end}}", // eqメソッドの引数の型がアンマッチで実行時にエラー
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != `template: template:1:5: executing "template" at <eq 1 "1">: error calling eq: incompatible types for comparison` {
		t.Fatal("failed test\n", err)
	}
}
