package cmd

import (
	"os"
	"testing"
)

func TestIncludeCmd(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
4,5
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1,col2
2,x
3,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1,col2",
		"2,3",
		"3,4",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestIncludeCmd_columnAnother(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
4,5
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1,col2
2,3
3,4
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"--column-another", "col2",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1,col2",
		"3,4",
		"4,5",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestIncludeCmd_duplicate(t *testing.T) {

	si := `col1,col2
1,2
2,3
1,x
3,4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
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
		"1,x",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestIncludeCmd_unmatch(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1
4
11
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1,col2",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestIncludeCmd_format(t *testing.T) {

	si := `col1	col2
1	2
2	3
3	4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1	col2
2	
3	
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"-o", fo.Name(),
		"--delim", `\t`,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo.Name())

	expect := joinRows(
		"col1	col2",
		"2	3",
		"3	4",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestIncludeCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fa := createTempFile(t, "")
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"-o", fo.Name(),
		"--delim", "\t\t",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestIncludeCmd_inputColumnNotFound(t *testing.T) {

	si := `col1,col2
1,2
2,3
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col3",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing col3 in the input CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestIncludeCmd_anthorColumnNotFound(t *testing.T) {

	si := `col1,col2
1,2
2,3
`
	fi := createTempFile(t, si)
	defer os.Remove(fi.Name())

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"--column-another", "col3",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the another CSV file: col3 is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestIncludeCmd_inputEmpty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name(),
		"-c", "col1",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the input CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestIncludeCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fa := createTempFile(t, "")
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name() + "____", // 存在しないファイル
		"-a", fa.Name(),
		"-c", "col1",
		"-o", fo.Name(),
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

func TestIncludeCmd_anotherFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi.Name())

	fa := createTempFile(t, "")
	defer os.Remove(fa.Name())

	fo := createTempFile(t, "")
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"include",
		"-i", fi.Name(),
		"-a", fa.Name() + "____", // 存在しないファイル
		"-c", "col1",
		"-o", fo.Name(),
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fa.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}
