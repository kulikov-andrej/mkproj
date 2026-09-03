package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateProject(t *testing.T) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"cpp",
		map[string]string{
			"CMakeLists.txt": "project(example)\n",
			"src/main.cpp":   "int main() {}\n",
			".gitignore":     "/build/\n",
		},
	)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	if err := env.app.run([]string{
		"-t",
		"cpp",
		target,
	}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(
		t,
		filepath.Join(
			target,
			"CMakeLists.txt",
		),
		"project(example)\n",
	)

	assertFileContent(
		t,
		filepath.Join(
			target,
			"src",
			"main.cpp",
		),
		"int main() {}\n",
	)

	assertFileContent(
		t,
		filepath.Join(
			target,
			".gitignore",
		),
		"/build/\n",
	)
}

func TestRejectNonEmptyTarget(t *testing.T) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"cpp",
		map[string]string{
			"main.cpp": "",
		},
	)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	if err := os.MkdirAll(
		target,
		0o755,
	); err != nil {
		t.Fatal(err)
	}

	existing := filepath.Join(
		target,
		"important.txt",
	)

	if err := os.WriteFile(
		existing,
		[]byte("keep me"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	err := env.app.run([]string{
		"-t",
		"cpp",
		target,
	})

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(
		err.Error(),
		"not empty",
	) {
		t.Fatalf("unexpected error: %v", err)
	}

	assertFileContent(
		t,
		existing,
		"keep me",
	)
}

func TestOpenProject(t *testing.T) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"cpp",
		map[string]string{
			"main.cpp": "",
		},
	)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	var openedPath string

	env.app.openProject = func(
		path string,
	) error {
		openedPath = path
		return nil
	}

	if err := env.app.run([]string{
		"-t",
		"cpp",
		target,
		"-o",
	}); err != nil {
		t.Fatal(err)
	}

	expected, err := filepath.Abs(target)
	if err != nil {
		t.Fatal(err)
	}

	if openedPath != filepath.Clean(expected) {
		t.Fatalf(
			"expected %q, got %q",
			expected,
			openedPath,
		)
	}
}
