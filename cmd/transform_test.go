package cmd

import (
	"os"
	"testing"
)

func TestTransformCmd(t *testing.T) {

	s := `ID,Name
1,Taro; Yamada
2,"Hanako, Sato"
`
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--out-delim", ";",
		"--out-quote", "'",
		"--out-sep", "|",
		"--out-allquote",
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

	expect := "'ID';'Name'|'1';'Taro; Yamada'|'2';'Hanako, Sato'|"
	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_custom(t *testing.T) {

	s := "ID/Name%1/Taro; Yamada%2/$Hanako, Sato$%"

	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--delim", "/",
		"--quote", "$",
		"--sep", "%",
		"--out-delim", ";",
		"--out-quote", "'",
		"--out-sep", "|",
		"--out-allquote",
		"--out-bom",
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

	expect := "\uFEFF'ID';'Name'|'1';'Taro; Yamada'|'2';'Hanako, Sato'|"
	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_fileNotFound(t *testing.T) {

	fi, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name() + "____", // 存在しないファイル名を指定
		"-o", fo.Name(),
		"--out-allquote",
	})

	err = rootCmd.Execute()
	if err == nil {
		t.Fatal("failed test\n", err)
	}

	pathErr := err.(*os.PathError)
	if pathErr.Path != fi.Name()+"____" || pathErr.Op != "open" {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_empty(t *testing.T) {

	s := ""

	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--out-allquote",
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

	if result != "" {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_delim(t *testing.T) {

	s := `ID	Name
1	Taro
2	"Hanako	Sato"
`
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--delim", `\t`,
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

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,Hanako	Sato\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_delim_multichar(t *testing.T) {

	s := ""
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--delim", `;;`,
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "flag delim should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_delim_multibyte(t *testing.T) {

	s := `ID　Name
1　Taro
2　"Hanako　Sato"
`
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--delim", `\u3000`, // マルチバイトとなる全角スペース
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

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,Hanako　Sato\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_delim_parseError(t *testing.T) {

	s := ""
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--delim", `\t"`,
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != `Could not parse value \t" of flag delim: invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_quote(t *testing.T) {

	s := `'ID','Name'
'1',Taro
2,'Hanako, Sato'
`
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--quote", "'",
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

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,\"Hanako, Sato\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_quote_multichar(t *testing.T) {

	s := ""
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--quote", "''",
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != "flag quote should be specified with a single character" {
		t.Fatal("failed test\n", err)
	}
}

func TestTransformCmd_sep(t *testing.T) {

	s := `ID,Name|1,Taro|2,"Hanako, Sato"`

	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--sep", "|",
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

	expect := "ID,Name\r\n" +
		"1,Taro\r\n" +
		"2,\"Hanako, Sato\"\r\n"

	if result != expect {
		t.Fatal("failed test\n", result)
	}
}

func TestTransformCmd_sep_parseError(t *testing.T) {

	s := ""
	fi, err := createTempFile(s)
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fi.Name())

	fo, err := createTempFile("")
	if err != nil {
		t.Fatal("failed test\n", err)
	}
	defer os.Remove(fo.Name())

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"transform",
		"-i", fi.Name(),
		"-o", fo.Name(),
		"--sep", `\r"`,
	})

	err = rootCmd.Execute()
	if err == nil || err.Error() != `Could not parse value \r" of flag sep: invalid syntax` {
		t.Fatal("failed test\n", err)
	}
}
