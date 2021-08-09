package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Test_execute(t *testing.T) {
	generatedFilePath := filepath.Join(os.TempDir(), "result.go")

	cmd := exec.Command(
		"go", "run", "./",
		"-output", generatedFilePath,
		"./test.sql",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("stdout", stdout.String())
	t.Log("stderr", stderr.String())

	b, err := os.ReadFile(generatedFilePath)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}
