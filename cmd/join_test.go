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
	f1 := createTempFile(t, s1)
	defer os.Remove(f1)

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	f2 := createTempFile(t, s2)
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1,
		"-2", f2,
		"-o", fo,
		"-c", "CompanyID",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID,CompanyName\r\n" +
		"1,Yamada,1,CompanyA\r\n" +
		"5,Ichikawa,1,CompanyA\r\n" +
		"2,\"Hanako, Sato\",3,会社C\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestJoinCmd_format(t *testing.T) {

	s1 := "ID;Name;CompanyID|1;Yamada;1|5;Ichikawa;1|2;'Hanako; Sato';3"
	f1 := createTempFile(t, s1)
	defer os.Remove(f1)

	s2 := "CompanyID;CompanyName|1;CompanyA|2;CompanyB|3;会社C|"
	f2 := createTempFile(t, s2)
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1,
		"-2", f2,
		"-o", fo,
		"-c", "CompanyID",
		"--delim", ";",
		"--quote", "'",
		"--sep", "|",
		"--allquote",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "'ID';'Name';'CompanyID';'CompanyName'|" +
		"'1';'Yamada';'1';'CompanyA'|" +
		"'5';'Ichikawa';'1';'CompanyA'|" +
		"'2';'Hanako; Sato';'3';'会社C'|"

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
	f1 := createTempFile(t, s1)
	defer os.Remove(f1)

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	f2 := createTempFile(t, s2)
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1,
		"-2", f2,
		"-o", fo,
		"-c", "CompanyID",
		"--usingfile",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID,CompanyName\r\n" +
		"1,Yamada,1,CompanyA\r\n" +
		"5,Ichikawa,1,CompanyA\r\n" +
		"2,\"Hanako, Sato\",3,会社C\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestJoinCmd_invalidFormat(t *testing.T) {

	f1 := createTempFile(t, "")
	defer os.Remove(f1)

	f2 := createTempFile(t, "")
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1,
		"-2", f2,
		"-o", fo,
		"-c", "CompanyID",
		"--delim", ";;",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestRunJoin_norecord(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,2
2,"Hanako, Sato",3
`
	f1 := createTempFile(t, s1)
	defer os.Remove(f1)

	s2 := `CompanyID,CompanyName
1,CompanyA
3,会社C
`
	f2 := createTempFile(t, s2)
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1,
		"-2", f2,
		"-o", fo,
		"-c", "CompanyID",
		"--norecord",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID,CompanyName\r\n" +
		"1,Yamada,1,CompanyA\r\n" +
		"5,Ichikawa,2,\r\n" +
		"2,\"Hanako, Sato\",3,会社C\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRunJoin_columnSecond(t *testing.T) {

	s1 := `ID,Name,CompanyID
1,Yamada,1
5,Ichikawa,2
2,"Hanako, Sato",3
`
	f1 := createTempFile(t, s1)
	defer os.Remove(f1)

	s2 := `ID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
4,会社D
`
	f2 := createTempFile(t, s2)
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"join",
		"-1", f1,
		"-2", f2,
		"-o", fo,
		"-c", "CompanyID",
		"--column-second", "ID",
	})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := readString(t, fo)

	expect := "ID,Name,CompanyID,CompanyName\r\n" +
		"1,Yamada,1,CompanyA\r\n" +
		"5,Ichikawa,2,CompanyB\r\n" +
		"2,\"Hanako, Sato\",3,会社C\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestRunJoin_firstFileNotFound(t *testing.T) {

	f1 := createTempFile(t, "")
	defer os.Remove(f1)

	f2 := createTempFile(t, "")
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	// 存在しないファイルを指定
	err := runJoin(csv.Format{}, f1+"___", f2, "CompanyID", fo, JoinOptions{})
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != f1+"___" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestRunJoin_secondFileNotFound(t *testing.T) {

	f1 := createTempFile(t, "")
	defer os.Remove(f1)

	f2 := createTempFile(t, "")
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	// 存在しないファイルを指定
	err := runJoin(csv.Format{}, f1, f2+"___", "CompanyID", fo, JoinOptions{})
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != f2+"___" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestRunJoin_outputFileNotFound(t *testing.T) {

	f1 := createTempFile(t, "")
	defer os.Remove(f1)

	f2 := createTempFile(t, "")
	defer os.Remove(f2)

	fo := createTempFile(t, "")
	defer os.Remove(fo)

	// 存在しないディレクトリのファイルを指定
	err := runJoin(csv.Format{}, f1, f2, "CompanyID", filepath.Join(fo, "___"), JoinOptions{})
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != filepath.Join(fo, "___") || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestJoin(t *testing.T) {

	s1 := `ID,Name
1,Yamada
5,Ichikawa
2,"Hanako, Sato"
`
	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := `ID,Height,Weight
1,171,50
2,160,60
5,152,50
`
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "ID", out, JoinOptions{})

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	out.Flush()
	result := b.String()

	expect := "ID,Name,Height,Weight\r\n" +
		"1,Yamada,171,50\r\n" +
		"5,Ichikawa,152,50\r\n" +
		"2,\"Hanako, Sato\",160,60\r\n"

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
	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := `ID,Height,Weight
5,152,50
`
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "ID", out, JoinOptions{})

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
	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := `ID,Height,Weight
5,152,50
`
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "ID", out, JoinOptions{noRecordNoError: true})

	if err != nil {
		t.Fatal("failed test\n", err)
	}

	out.Flush()
	result := b.String()

	expect := "ID,Name,Height,Weight\r\n" +
		"1,Yamada,,\r\n" +
		"5,Ichikawa,152,50\r\n" +
		"2,\"Hanako, Sato\",,\r\n"

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
	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "CompanyID", out, JoinOptions{})
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
	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := `ID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "CompanyID", out, JoinOptions{})
	if err == nil || err.Error() != "failed to read the second CSV file: CompanyID is not found" {
		t.Fatal("failed test\n", err)
	}
}

func TestJoin_firstFileEmpty(t *testing.T) {

	s1 := ""

	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := `CompanyID,CompanyName
1,CompanyA
2,CompanyB
3,会社C
`
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "CompanyID", out, JoinOptions{})
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
	r1 := csv.NewCsvReader(strings.NewReader(s1), csv.Format{})

	s2 := ""
	r2 := csv.NewCsvReader(strings.NewReader(s2), csv.Format{})

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	out := csv.NewCsvWriter(w, csv.Format{})

	err := join(r1, r2, "CompanyID", out, JoinOptions{})
	if err == nil || err.Error() != "failed to read the second CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
