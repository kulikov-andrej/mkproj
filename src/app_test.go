package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testApp struct {
	app          *app
	templateRoot string
	workRoot     string
	stdout       *bytes.Buffer
	stderr       *bytes.Buffer
}

func newTestApp(t *testing.T) *testApp {
	t.Helper()

	root := t.TempDir()

	templateRoot := filepath.Join(root, "templates")
	workRoot := filepath.Join(root, "work")

	if err := os.MkdirAll(templateRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(workRoot, 0o755); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	a := newApp(
		templateRoot,
		strings.NewReader(""),
		&stdout,
		&stderr,
	)

	return &testApp{
		app:          a,
		templateRoot: templateRoot,
		workRoot:     workRoot,
		stdout:       &stdout,
		stderr:       &stderr,
	}
}

func (env *testApp) createTemplate(
	t *testing.T,
	name string,
	files map[string]string,
) string {
	t.Helper()

	root := filepath.Join(
		env.templateRoot,
		name,
	)

	for path, content := range files {
		fullPath := filepath.Join(root, path)

		if err := os.MkdirAll(
			filepath.Dir(fullPath),
			0o755,
		); err != nil {
			t.Fatal(err)
		}

		if err := os.WriteFile(
			fullPath,
			[]byte(content),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	return root
}
