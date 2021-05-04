package cmd

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/onozaty/csvt/csv"
)

func TestJoinCmd(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	f1, err := createTempFile(s1)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	f2, err := createTempFile(s2)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1.Name(),
		"-2", f2.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

	expect := `ID,Name,CompanyID,CompanyName
1,Yamada,1,CompanyA
5,Ichikawa,1,CompanyA
2,"Hanako, Sato",3,会社C
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestJoinCmd_usingfile(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	f1, err := createTempFile(s1)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	f2, err := createTempFile(s2)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1.Name(),
		"-2", f2.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"--usingfile",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

	expect := `ID,Name,CompanyID,CompanyName
1,Yamada,1,CompanyA
5,Ichikawa,1,CompanyA
2,"Hanako, Sato",3,会社C
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRunJoin_norecord(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,2
2,"Hanako, Sato",3
`
	f1, err := createTempFile(s1)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	s2 := `CompanyID,CompanyName
1,CompanyA
3,会社C
`
	f2, err := createTempFile(s2)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1.Name(),
		"-2", f2.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"--norecord",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

	expect := `ID,Name,CompanyID,CompanyName
1,Yamada,1,CompanyA
5,Ichikawa,2,
2,"Hanako, Sato",3,会社C
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRunJoin_column2(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,2
2,"Hanako, Sato",3
`
	f1, err := createTempFile(s1)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	s2 := `ID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
4,会社D
`
	f2, err := createTempFile(s2)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1.Name(),
		"-2", f2.Name(),
		"-o", fo.Name(),
		"-c", "CompanyID",
		"--column2", "ID",
	})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	bo, err := os.ReadFile(fo.Name())
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := string(bo)

	expect := `ID,Name,CompanyID,CompanyName
1,Yamada,1,CompanyA
5,Ichikawa,2,CompanyB
2,"Hanako, Sato",3,会社C
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRunJoin_firstFileNotFound(t *testing.T) {

	f1, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	f2, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	// 存在しないファイルを指定
	err = runJoin(f1.Name()+"___", f2.Name(), "CompanyID", fo.Name(), JoinOptions{})
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != f1.Name()+"___" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestRunJoin_secondFileNotFound(t *testing.T) {

	f1, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	f2, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	// 存在しないファイルを指定
	err = runJoin(f1.Name(), f2.Name()+"___", "CompanyID", fo.Name(), JoinOptions{})
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != f2.Name()+"___" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestRunJoin_outputFileNotFound(t *testing.T) {

	f1, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f1.Name())

	f2, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(f2.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	// 存在しないディレクトリのファイルを指定
	err = runJoin(f1.Name(), f2.Name(), "CompanyID", filepath.Join(fo.Name(), "___"), JoinOptions{})
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != filepath.Join(fo.Name(), "___") || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestJoin(t *testing.T) {

	s1 := `ID,Name
1,Yamada
5,Ichikawa
2,"Hanako, Sato"
`
	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := `ID,Height,Weight
1,171,50
2,160,60
5,152,50
`
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "ID", out, JoinOptions{})

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	out.Flush()
	result := b.String()

	expect := `ID,Name,Height,Weight
1,Yamada,171,50
5,Ichikawa,152,50
2,"Hanako, Sato",160,60
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestJoin_rightNoneError(t *testing.T) {

	s1 := `ID,Name
1,Yamada
5,Ichikawa
2,"Hanako, Sato"
`
	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := `ID,Height,Weight
5,152,50
`
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "ID", out, JoinOptions{})

	if err == nil || err.Error() != "1 was not found in the second CSV file\nif you don't want to raise an error, use the 'norecord' option" {
		t.Fatal("failed test\n", err)
	}

}

func TestJoin_rightNoneNoError(t *testing.T) {

	s1 := `ID,Name
1,Yamada
5,Ichikawa
2,"Hanako, Sato"
`
	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := `ID,Height,Weight
5,152,50
`
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "ID", out, JoinOptions{noRecordNoError: true})

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	out.Flush()
	result := b.String()

	expect := `ID,Name,Height,Weight
1,Yamada,,
5,Ichikawa,152,50
2,"Hanako, Sato",,
`

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestJoin_firstFileJoinColumnNotFound(t *testing.T) {

	s1 := `ID,Name,CID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "CompanyID", out, JoinOptions{})
	if err == nil || err.Error() != "missing CompanyID in the first CSV file" {
		t.Fatal("failed test\n", err)
	}
}

func TestJoin_secondFileJoinColumnNotFound(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := `ID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "CompanyID", out, JoinOptions{})
	if err == nil || err.Error() != "failed to read the second CSV file: CompanyID is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestJoin_firstFileEmpty(t *testing.T) {

	s1 := ""

	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "CompanyID", out, JoinOptions{})
	if err == nil || err.Error() != "failed to read the first CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func TestJoin_secondFileEmpty(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,1
2,"Hanako, Sato",3
`
	r1, err := csv.NewCsvReader(strings.NewReader(s1))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	s2 := ""
	r2, err := csv.NewCsvReader(strings.NewReader(s2))
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w)

	err = join(r1, r2, "CompanyID", out, JoinOptions{})
	if err == nil || err.Error() != "failed to read the second CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}

func createTempFile(content string) (*os.File, error) {

	tempFile, err := os.CreateTemp("", "csv")
	if err != nil {
		return nil, err
	}

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}
