package cmd

import (
	"os"
	"testing"
)

func TestSliceCmd(t *testing.T) {

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

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "2",
		"-e", "4",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"2,b",
		"3,c",
		"4,d",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSliceCmd_startOnly(t *testing.T) {

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

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"3,c",
		"4,d",
		"5,e",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSliceCmd_endOnly(t *testing.T) {

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

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-e", "3",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,a",
		"2,b",
		"3,c",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSliceCmd_sameStartEnd(t *testing.T) {

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

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "1",
		"-e", "1",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,a",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSliceCmd_format(t *testing.T) {

	s := "col1,col2;1,a;2,b;"

	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "2",
		"--sep", ";",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "col1,col2;2,b;"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestSliceCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "2",
		"--quote", ";;",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag quote should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestSliceCmd_start0(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "0",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "start must be greater than or equal to 1" {
		t.Fatal("failed test\n", err)
	}
}

func TestSliceCmd_endStart(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "10",
		"-e", "9",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "end must be greater than or equal to start" {
		t.Fatal("failed test\n", err)
	}
}

func TestSliceCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi,
		"-o", fo,
		"-s", "1",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the input CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestSliceCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"slice",
		"-i", fi + "____", // 存在しないファイル
		"-o", fo,
		"-s", "1",
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
