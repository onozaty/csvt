package cmd

import (
	"os"
	"testing"
)

func createTempFile(t *testing.T, content string) *os.File {

	tempFile, err := os.CreateTemp("", "csv")
	if err != nil {
		t.Fatal("craete file failed\n", err)
	}

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	return tempFile
}

func readString(t *testing.T, name string) string {

	bo, err := os.ReadFile(name)
	if err != nil {
		t.Fatal("read failed\n", err)
	}

	return string(bo)
}
