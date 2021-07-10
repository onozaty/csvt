package cmd

import (
	"os"
	"testing"
)

func TestUniqueCmd(t *testing.T) {

	s := `col1,col2,col3
1,2,3
2,2,2
1,1,1
1,2,3
3,2,1
3,3,3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,2,3",
		"2,2,2",
		"3,2,1",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestUniqueCmd_format(t *testing.T) {

	s := `col1	col2	col3
1	2	3
1	2	3
2	2	3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col1",
		"-o", fo,
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1\tcol2\tcol3",
		"1\t2\t3",
		"2\t2\t3",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestUniqueCmd_multiColumn(t *testing.T) {

	s := `col1,col2,col3
1,11,3
11,1,2
1,11,1
1,2,3
3,2,1
2,3,3
2,3,1
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col1",
		"-c", "col2",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2,col3",
		"1,11,3",
		"11,1,2",
		"1,2,3",
		"3,2,1",
		"2,3,3",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestUniqueCmd_columnNotFound(t *testing.T) {

	s := `col1,col2,col3
1,2,3
`
	fi := createTempFile(t, s)
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col4",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing col4 in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestUniqueCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col1",
		"-o", fo,
		"--delim", "zz",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestUniqueCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestUniqueCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi + "____", // 存在しないファイル
		"-c", "col1",
		"-o", fo,
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

func TestUniqueCmd_outputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"unique",
		"-i", fi,
		"-c", "col1",
		"-o", fo + "/___", // 存在しないディレクトリ
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fo+"/___" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}
