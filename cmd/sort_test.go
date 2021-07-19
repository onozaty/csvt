package cmd

import (
	"io"
	"os"
	"testing"
)

func TestSortCmd(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"2,a",
		"1,b",
		"4,c",
		"3,d",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,b",
		"2,a",
		"3,d",
		"4,c",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_format(t *testing.T) {

	s := joinRows(
		"col1\tcol2",
		"2\ta",
		"1\tb",
		"4\tc",
		"3\td",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1\tcol2",
		"1\tb",
		"2\ta",
		"3\td",
		"4\tc",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}
func TestSortCmd_multiColumn(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"2,b",
		"1,b",
		"2,a",
		"3,a",
		"1,c",
		"1,a",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"-c", "col2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,a",
		"1,b",
		"1,c",
		"2,a",
		"2,b",
		"3,a",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_number(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"100,b",
		"1,b",
		"11,a",
		"2,a",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"--number",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,b",
		"2,a",
		"11,a",
		"100,b",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_desc(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"2,a",
		"1,b",
		"4,c",
		"3,d",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"--desc",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"4,c",
		"3,d",
		"2,a",
		"1,b",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_multiColumn_number_desc(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,100",
		"2,10",
		"11,10",
		"4,2",
		"5,1",
		"10,10",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col2",
		"-c", "col1",
		"--number",
		"--desc",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,100",
		"11,10",
		"10,10",
		"2,10",
		"4,2",
		"5,1",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_usingfile(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"2,a",
		"1,b",
		"4,c",
		"3,d",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"--usingfile",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,b",
		"2,a",
		"3,d",
		"4,c",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_multiColumn_number_desc_usingfile(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,100",
		"2,10",
		"11,10",
		"4,2",
		"5,1",
		"10,10",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col2",
		"-c", "col1",
		"--number",
		"--desc",
		"--usingfile",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,100",
		"11,10",
		"10,10",
		"2,10",
		"4,2",
		"5,1",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSortCmd_invalidFormat(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,1",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
		"--delim", "xx",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestSortCmd_columnNotFound(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,1",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "col3 is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestSortCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
	})

	err := rootCmd.Execute()
	if err != io.EOF {
		t.Fatal("failed test\n", err)
	}
}

func TestSortCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"sort",
		"-i", fi + "____", // 存在しないファイル
		"-o", fo,
		"-c", "col1",
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
