package cmd

import (
	"os"
	"testing"
)

func TestRenameCmd(t *testing.T) {

	s := `A,B,C,D
1,x,a,_
2,y,b,_
3,z,c,_
`

	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "B",
		"-a", "B-before",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := "A,B-before,C,D\r\n" +
		"1,x,a,_\r\n" +
		"2,y,b,_\r\n" +
		"3,z,c,_\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRenameCmd_format(t *testing.T) {

	s := `A	B
1	x
`

	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "B",
		"-a", "B-before",
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := "A\tB-before\r\n" +
		"1\tx\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRenameCmd_columns(t *testing.T) {

	s := `A,B,C,D
1,x,a,_
2,y,b,_
3,z,c,_
`

	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "C",
		"-a", "A",
		"-c", "A",
		"-a", "C",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := "C,B,A,D\r\n" +
		"1,x,a,_\r\n" +
		"2,y,b,_\r\n" +
		"3,z,c,_\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRenameCmd_fileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name() + "____", // 存在しないファイル名を指定
		"-o", fo.Name(),
		"-c", "before",
		"-a", "after",
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

func TestRenameCmd_columnNotFound(t *testing.T) {

	s := `A,B,C,D
1,x,a,_
2,y,b,_
3,z,c,_
`
	fi := createTempFile(t, s)
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "a", // 存在しないカラム
		"-a", "after",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing a in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestRenameCmd_empty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "A",
		"-a", "a",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestRenameCmd_column_unmatched(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "A",
		"-c", "B",
		"-a", "a",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "the number of columns before and after the renaming is unmatched" {
		t.Fatal("failed test\n", err)
	}
}

func TestRenameCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"rename",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"-c", "A",
		"-a", "a",
		"--delim", "__",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}
