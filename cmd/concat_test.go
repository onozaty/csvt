package cmd

import (
	"os"
	"testing"
)

func TestConcatCmd(t *testing.T) {

	s1 := `col1,col2
1,2
2,3
3,4
`
	fi1 := createTempFile(t, s1)
	defer os.Remove(fi1.Name())

	s2 := `col1,col2
2,x
3,y
`
	fi2 := createTempFile(t, s2)
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1,col2",
		"1,2",
		"2,3",
		"3,4",
		"2,x",
		"3,y",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestConcatCmd_swap(t *testing.T) {

	s1 := `col1,col2,col3
1,x2,x3
2,y2,y3
`
	fi1 := createTempFile(t, s1)
	defer os.Remove(fi1.Name())

	s2 := `col2,col3,col1
a2,a3,3
b2,b3,4
`
	fi2 := createTempFile(t, s2)
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1,col2,col3",
		"1,x2,x3",
		"2,y2,y3",
		"3,a2,a3",
		"4,b2,b3",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestConcatCmd_format(t *testing.T) {

	s1 := `col1	col2
1	2
`
	fi1 := createTempFile(t, s1)
	defer os.Remove(fi1.Name())

	s2 := `col1	col2
a	b
`
	fi2 := createTempFile(t, s2)
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1\tcol2",
		"1\t2",
		"a\tb",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestConcatCmd_invalidFormat(t *testing.T) {

	fi1 := createTempFile(t, "")
	defer os.Remove(fi1.Name())

	fi2 := createTempFile(t, "")
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
		"--delim", `\t\t`,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_columnCountUnmatch(t *testing.T) {

	s1 := `col1,col2
1,2
`
	fi1 := createTempFile(t, s1)
	defer os.Remove(fi1.Name())

	s2 := `col1,col2,col3
2,x,y
`
	fi2 := createTempFile(t, s2)
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "number of columns does not match" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_columnNotFound(t *testing.T) {

	s1 := `col1,col2
1,2
`
	fi1 := createTempFile(t, s1)
	defer os.Remove(fi1.Name())

	s2 := `col1,col3
2,x
`
	fi2 := createTempFile(t, s2)
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "no column corresponding to the second CSV file: missing col2 in the CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_firstEmpty(t *testing.T) {

	fi1 := createTempFile(t, "")
	defer os.Remove(fi1.Name())

	s2 := `col1,col2
2,x
`
	fi2 := createTempFile(t, s2)
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the first CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_secondEmpty(t *testing.T) {

	s1 := `col1,col2
2,x
`

	fi1 := createTempFile(t, s1)
	defer os.Remove(fi1.Name())

	fi2 := createTempFile(t, "")
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the second CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_firstFileNotFound(t *testing.T) {

	fi1 := createTempFile(t, "")
	defer os.Remove(fi1.Name())

	fi2 := createTempFile(t, "")
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name() + "____", // 存在しないファイル
		"-2", fi2.Name(),
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fi1.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_secondFileNotFound(t *testing.T) {

	fi1 := createTempFile(t, "")
	defer os.Remove(fi1.Name())

	fi2 := createTempFile(t, "")
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name() + "____", // 存在しないファイル
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fi2.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestConcatCmd_outputFileNotFound(t *testing.T) {

	fi1 := createTempFile(t, "")
	defer os.Remove(fi1.Name())

	fi2 := createTempFile(t, "")
	defer os.Remove(fi2.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"concat",
		"-1", fi1.Name(),
		"-2", fi2.Name(),
		"-o", fo.Name() + "/aa", // 存在しないフォルダ
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fo.Name()+"/aa" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}