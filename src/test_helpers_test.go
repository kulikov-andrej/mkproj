package main

import (
	"os"
	"testing"
)

func assertFileContent(
	t *testing.T,
	path string,
	expected string,
) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %q: %v", path, err)
	}

	if string(content) != expected {
		t.Fatalf(
			"%q: expected %q, got %q",
			path,
			expected,
			content,
		)
	}
}

func assertNotExists(
	t *testing.T,
	path string,
) {
	t.Helper()

	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		t.Fatalf(
			"expected %q not to exist",
			path,
		)
	}
}
