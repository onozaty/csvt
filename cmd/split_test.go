package cmd

import (
	"os"
	"reflect"
	"testing"
)

func TestSplitCmd(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
		"4,d",
		"5,e",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"output-1.csv": []byte(joinRows(
			"col1,col2",
			"1,a",
			"2,b")),
		"output-2.csv": []byte(joinRows(
			"col1,col2",
			"3,c",
			"4,d")),
		"output-3.csv": []byte(joinRows(
			"col1,col2",
			"5,e")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_maxRows_1(t *testing.T) {

	s := joinRows(
		"col1",
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "1",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"output-1.csv":  []byte(joinRows("col1", "1")),
		"output-2.csv":  []byte(joinRows("col1", "2")),
		"output-3.csv":  []byte(joinRows("col1", "3")),
		"output-4.csv":  []byte(joinRows("col1", "4")),
		"output-5.csv":  []byte(joinRows("col1", "5")),
		"output-6.csv":  []byte(joinRows("col1", "6")),
		"output-7.csv":  []byte(joinRows("col1", "7")),
		"output-8.csv":  []byte(joinRows("col1", "8")),
		"output-9.csv":  []byte(joinRows("col1", "9")),
		"output-10.csv": []byte(joinRows("col1", "10")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_maxRows_over_rows(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "10",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"output-1.csv": []byte(joinRows(
			"col1,col2",
			"1,a",
			"2,b",
			"3,c")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_maxRows_equal_rows(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"output-1.csv": []byte(joinRows(
			"col1,col2",
			"1,a",
			"2,b",
			"3,c")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_format(t *testing.T) {

	s := "col1\tcol2;1\ta;2\tb;3\tc;4\td;"

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "2",
		"--delim", "\t",
		"--sep", ";",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"output-1.csv": []byte("col1\tcol2;1\ta;2\tb;"),
		"output-2.csv": []byte("col1\tcol2;3\tc;4\td;"),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "2",
		"--delim", "aa",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestSplitCmd_maxRows_0(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "0",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "rows must be greater than or equal to 1" {
		t.Fatal("failed test\n", err)
	}
}

func TestSplitCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi + "____", // 存在しないファイル
		"-o", d + "/output.csv",
		"-r", "1",
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

func TestSplitCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the input CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestSplitCmd_createParent(t *testing.T) {

	s := joinRows(
		"col1",
		"1",
		"2",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/parent/output.csv",
		"-r", "1",
	})

	// 親ディレクトリだけは作られる
	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d+"/parent")

	expect := map[string][]byte{
		"output-1.csv": []byte(joinRows(
			"col1",
			"1")),
		"output-2.csv": []byte(joinRows(
			"col1",
			"2")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_canotCreateParent(t *testing.T) {

	s := joinRows(
		"col1",
		"1",
		"2",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/parent1/parent2/output.csv",
		"-r", "1",
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Op != "mkdir" {
		t.Fatal("failed test\n", err)
	}
}

func TestSplitCmd_headerOnly(t *testing.T) {

	s := "col1,col2"

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/output.csv",
		"-r", "2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"output-1.csv": []byte(joinRows("col1,col2")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_outputBase_paramZeroPadding(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"\"1\",a",
		"2,\"b\"",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/%04d.csv",
		"-r", "2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"0001.csv": []byte(joinRows(
			"col1,col2",
			"1,a",
			"2,b")),
		"0002.csv": []byte(joinRows(
			"col1,col2",
			"3,c")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_outputBase_param(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"\"1\",a",
		"2,\"b\"",
		"3,c",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/out%d.csv",
		"-r", "2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"out1.csv": []byte(joinRows(
			"col1,col2",
			"1,a",
			"2,b")),
		"out2.csv": []byte(joinRows(
			"col1,col2",
			"3,c")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}

func TestSplitCmd_outputBase_nonExtension(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"\"1\",a",
		"2,\"b\"",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	d := createTempDir(t)
	defer os.RemoveAll(d)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"split",
		"-i", fi,
		"-o", d + "/a",
		"-r", "2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readDir(t, d)

	expect := map[string][]byte{
		"a-1": []byte(joinRows(
			"col1,col2",
			"1,a",
			"2,b")),
	}

	if !reflect.DeepEqual(result, expect) {
		t.Fatal("failed test\n", result)
	}
}
