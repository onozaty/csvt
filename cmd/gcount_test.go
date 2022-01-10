package cmd

import (
	"os"
	"testing"
)

func TestGcountCmd(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,B",
		"2,A",
		"3,a",
		"4,A",
		"5,C",
		"6,C",
		"7,A",
		"8,AA",
		"9,B",
		"10,",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"gcount",
		"-i", fi,
		"-o", fo,
		"-c", "col2",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col2,COUNT",
		",1",
		"A,3",
		"AA,1",
		"B,2",
		"C,2",
		"a,1",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestGcountCmd_countColumn(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,A",
		"2,A",
		"3,A",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"gcount",
		"-i", fi,
		"-o", fo,
		"-c", "col2",
		"--count-column", "count",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col2,count",
		"A,3",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestGcountCmd_format(t *testing.T) {

	s := joinRows(
		"col1\tcol2",
		"1\ta",
		"2\tb",
		"3\ta",
		"4\tb",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"gcount",
		"-i", fi,
		"-o", fo,
		"-c", "col2",
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col2\tCOUNT",
		"a\t2",
		"b\t2",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestGcountCmd_invalidFormat(t *testing.T) {

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
		"gcount",
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

func TestGcountCmd_columnNotFound(t *testing.T) {

	s := joinRows(
		"col1,col2",
		"1,A",
		"2,A",
		"3,A",
	)

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"gcount",
		"-i", fi,
		"-o", fo,
		"-c", "col3",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing col3 in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestGcountCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"gcount",
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

func TestGcountCmd_inputFileEmpty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"gcount",
		"-i", fi,
		"-o", fo,
		"-c", "col1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
