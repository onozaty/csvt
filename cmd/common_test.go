package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTempFile(t *testing.T, content string) string {

	tempFile, err := os.CreateTemp("", "csv")
	if err != nil {
		t.Fatal("craete file failed\n", err)
	}

	_, err = tempFile.Write([]byte(content))
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	err = tempFile.Close()
	if err != nil {
		t.Fatal("write file failed\n", err)
	}

	return tempFile.Name()
}

func createTempDir(t *testing.T) string {

	tempDir, err := os.MkdirTemp("", "csvt")
	if err != nil {
		t.Fatal("craete dir failed\n", err)
	}

	return tempDir
}

func readBytes(t *testing.T, name string) []byte {

	bo, err := os.ReadFile(name)
	if err != nil {
		t.Fatal("read failed\n", err)
	}

	return bo
}

func readString(t *testing.T, name string) string {

	bo := readBytes(t, name)
	return string(bo)
}

func joinRows(rows ...string) string {
	return strings.Join(rows, "\r\n") + "\r\n"
}

func readDir(t *testing.T, dir string) map[string][]byte {

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal("read dir failed\n", err)
	}

	contentsMap := map[string][]byte{}
	for _, f := range files {
		b := readBytes(t, filepath.Join(dir, f.Name()))
		contentsMap[f.Name()] = b
	}

	return contentsMap
}
