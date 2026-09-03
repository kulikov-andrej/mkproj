package main

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplates(t *testing.T) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"python",
		map[string]string{
			"main.py": "",
		},
	)

	env.createTemplate(
		t,
		"cpp",
		map[string]string{
			"main.cpp": "",
		},
	)

	if err := env.app.run(
		[]string{"--list"},
	); err != nil {
		t.Fatal(err)
	}

	if got, want := env.stdout.String(), "cpp\npython\n"; got != want {
		t.Fatalf(
			"expected %q, got %q",
			want,
			got,
		)
	}
}

func TestMissingTemplate(t *testing.T) {
	env := newTestApp(t)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	err := env.app.run([]string{
		"-t",
		"missing",
		target,
	})

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(
		err.Error(),
		"not found",
	) {
		t.Fatalf("unexpected error: %v", err)
	}

	assertNotExists(t, target)
}
