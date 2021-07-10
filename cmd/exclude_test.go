package cmd

import (
	"os"
	"testing"
)

func TestExcludeCmd(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
4,5
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1,col2
2,x
3,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,2",
		"4,5",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestExcludeCmd_columnAnother(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
4,5
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1,col2
2,3
3,4
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"--column-another", "col2",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,2",
		"2,3",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestExcludeCmd_duplicate(t *testing.T) {

	si := `col1,col2
1,2
2,3
1,x
3,4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"2,3",
		"3,4",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestExcludeCmd_match_none(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1
4
11
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
		"1,2",
		"2,3",
		"3,4",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestExcludeCmd_match_all(t *testing.T) {

	si := `col1,col2
1,2
2,3
3,4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1
4
3
2
1
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := joinRows(
		"col1,col2",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestExcludeCmd_format(t *testing.T) {

	si := `col1	col2
1	2
2	3
3	4
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1	col2
2	
3	
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
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
		"col1	col2",
		"1	2",
	)

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestExcludeCmd_invalidFormat(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fa := createTempFile(t, "")
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"-o", fo,
		"--delim", "\t\t",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestExcludeCmd_inputColumnNotFound(t *testing.T) {

	si := `col1,col2
1,2
2,3
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col3",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "missing col3 in the input CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestExcludeCmd_anthorColumnNotFound(t *testing.T) {

	si := `col1,col2
1,2
2,3
`
	fi := createTempFile(t, si)
	defer os.Remove(fi)

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"--column-another", "col3",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the another CSV file: col3 is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestExcludeCmd_inputEmpty(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	sa := `col1,col2
1,x
1,y
`
	fa := createTempFile(t, sa)
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa,
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the input CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestExcludeCmd_inputFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fa := createTempFile(t, "")
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi + "____", // 存在しないファイル
		"-a", fa,
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

func TestExcludeCmd_anotherFileNotFound(t *testing.T) {

	fi := createTempFile(t, "")
	defer os.Remove(fi)

	fa := createTempFile(t, "")
	defer os.Remove(fa)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"exclude",
		"-i", fi,
		"-a", fa + "____", // 存在しないファイル
		"-c", "col1",
		"-o", fo,
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fa+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}
