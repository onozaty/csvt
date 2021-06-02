package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCmd(t *testing.T) {

	rootCmd := newRootCmd()
	rootCmd.SetArgs([]string{
		"version",
	})

	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)

	err := rootCmd.Execute()
	if err != nil {
		t.Fatal("failed test\n", err)
	}

	result := buf.String()

	if !strings.HasPrefix(result, "Version: dev\nRevision: dev\nOS: ") {
		t.Fatal("failed test\n", result)
	}
}
