package project

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/data/templates"
)

func TestCreate(t *testing.T) {
	root := t.TempDir()

	templatePath := filepath.Join(root, "template")
	target := filepath.Join(root, "hello")

	if err := os.MkdirAll(templatePath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		filepath.Join(templatePath, "main.txt"),
		[]byte("hello"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	got, err := Create(
		templates.Template{
			Name: "example",
			Path: templatePath,
		},
		target,
	)
	if err != nil {
		t.Fatal(err)
	}

	expectedPath, err := filepath.Abs(target)
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != "hello" {
		t.Fatalf(
			"expected project name %q, got %q",
			"hello",
			got.Name,
		)
	}

	if got.Path != expectedPath {
		t.Fatalf(
			"expected project path %q, got %q",
			expectedPath,
			got.Path,
		)
	}

	content, err := os.ReadFile(
		filepath.Join(target, "main.txt"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != "hello" {
		t.Fatalf(
			"expected %q, got %q",
			"hello",
			string(content),
		)
	}
}
func TestCreateRejectsNonEmptyTarget(t *testing.T) {
	root := t.TempDir()

	templatePath := filepath.Join(root, "template")
	target := filepath.Join(root, "target")

	if err := os.MkdirAll(templatePath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(target, 0o755); err != nil {
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

	_, err := Create(
		templates.Template{
			Name: "example",
			Path: templatePath,
		},
		target,
	)

	if err == nil {
		t.Fatal("expected an error")
	}

	content, readErr := os.ReadFile(existing)
	if readErr != nil {
		t.Fatal(readErr)
	}

	if string(content) != "keep me" {
		t.Fatal("existing target content was modified")
	}
}
func TestCreateRejectsFileTarget(t *testing.T) {
	root := t.TempDir()

	templatePath := filepath.Join(root, "template")
	target := filepath.Join(root, "target")

	if err := os.MkdirAll(templatePath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		target,
		[]byte("existing"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	_, err := Create(
		templates.Template{
			Name: "example",
			Path: templatePath,
		},
		target,
	)

	if err == nil {
		t.Fatal("expected an error")
	}
}
func TestCleanupMetadata(t *testing.T) {
	root := t.TempDir()

	metadata := filepath.Join(
		root,
		".mkproj",
	)

	if err := os.MkdirAll(metadata, 0o755); err != nil {
		t.Fatal(err)
	}

	projectFile := filepath.Join(
		root,
		"main.txt",
	)

	if err := os.WriteFile(
		projectFile,
		[]byte("keep"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	err := CleanupMetadata(Project{
		Name: "example",
		Path: root,
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(metadata); !os.IsNotExist(err) {
		t.Fatalf(
			"expected metadata to be removed, got %v",
			err,
		)
	}

	if _, err := os.Stat(projectFile); err != nil {
		t.Fatalf(
			"project file should remain: %v",
			err,
		)
	}
}
