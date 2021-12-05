package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestHeadCmd(t *testing.T) {

	s := `ID,Name,Company ID
1,Yamada,1
2,Ichikawa,
3,"Hanako, Sato",3
4,Otani,3
5,,2
6,Suzuki,1
7,Jane,2
8,Michel,3
9,1,3
10,Ken,3
11,Suzuki,2
12,Z,
`
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"head",
		"-i", f,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	except := `+----+--------------+------------+
| ID | Name         | Company ID |
+----+--------------+------------+
| 1  | Yamada       | 1          |
| 2  | Ichikawa     |            |
| 3  | Hanako, Sato | 3          |
| 4  | Otani        | 3          |
| 5  |              | 2          |
| 6  | Suzuki       | 1          |
| 7  | Jane         | 2          |
| 8  | Michel       | 3          |
| 9  | 1            | 3          |
| 10 | Ken          | 3          |
+----+--------------+------------+
`
	if result != except {
		t.Fatal("failed test\n", result)
	}
}

func TestHeadCmd_number(t *testing.T) {

	s := `ID,Name,Company ID
1,Yamada,1
2,Ichikawa,
3,"Hanako, Sato",3
4,Otani,3
5,,2
`
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"head",
		"-i", f,
		"-n", "3",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	except := `+----+--------------+------------+
| ID | Name         | Company ID |
+----+--------------+------------+
| 1  | Yamada       | 1          |
| 2  | Ichikawa     |            |
| 3  | Hanako, Sato | 3          |
+----+--------------+------------+
`
	if result != except {
		t.Fatal("failed test\n", result)
	}
}

func TestHeadCmd_less(t *testing.T) {

	s := `ID,Name,Company ID
1,Yamada,1
2,Ichikawa,
`
	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"head",
		"-i", f,
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	except := `+----+----------+------------+
| ID | Name     | Company ID |
+----+----------+------------+
| 1  | Yamada   | 1          |
| 2  | Ichikawa |            |
+----+----------+------------+
`
	if result != except {
		t.Fatal("failed test\n", result)
	}
}

func TestHeadCmd_format(t *testing.T) {

	s := "col1,col2;1,a;2,b;"

	f := createTempFile(t, s)
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"head",
		"-i", f,
		"--sep", ";",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	except := `+------+------+
| col1 | col2 |
+------+------+
| 1    | a    |
| 2    | b    |
+------+------+
`
	if result != except {
		t.Fatal("failed test\n", result)
	}
}

func TestHeadCmd_invalidNumber(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"head",
		"-i", f,
		"-n", "0",
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "number must be greater than or equal to 1" {
		t.Fatal("failed test\n", err)
	}
}

func TestHeadCmd_empty(t *testing.T) {

	f := createTempFile(t, "")
	defer os.Remove(f)

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"head",
		"-i", f,
	})

	err := rootCmd.Execute()
	if err == nil || err.Error() != "failed to read the input CSV file: EOF" {
		t.Fatal("failed test\n", err)
	}
}
