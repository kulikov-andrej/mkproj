package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func TestPrepareTargetNonexistent(t *testing.T) {
	target := filepath.Join(t.TempDir(), "new-project")

	got, err := prepareTarget(target)
	if err != nil {
		t.Fatalf("prepareTarget() error = %v", err)
	}

	want, err := filepath.Abs(target)
	if err != nil {
		t.Fatal(err)
	}
	want = filepath.Clean(want)

	if got != want {
		t.Fatalf("prepareTarget() = %q, want %q", got, want)
	}
}

func TestPrepareTargetExistingEmptyDirectory(t *testing.T) {
	target := t.TempDir()

	got, err := prepareTarget(target)
	if err != nil {
		t.Fatalf("prepareTarget() error = %v", err)
	}

	want, err := filepath.Abs(target)
	if err != nil {
		t.Fatal(err)
	}
	if got != filepath.Clean(want) {
		t.Fatalf("prepareTarget() = %q, want %q", got, filepath.Clean(want))
	}
}

func TestPrepareTargetErrors(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := prepareTarget("")
		if err == nil || err.Error() != "target path cannot be empty" {
			t.Fatalf("prepareTarget() error = %v", err)
		}
	})

	t.Run("file", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "target")
		if err := os.WriteFile(path, []byte("file"), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := prepareTarget(path)
		if err == nil || !strings.Contains(err.Error(), "target is not a directory") {
			t.Fatalf("prepareTarget() error = %v", err)
		}
	})

	t.Run("non-empty directory", func(t *testing.T) {
		path := t.TempDir()
		if err := os.WriteFile(filepath.Join(path, "file.txt"), []byte("data"), 0o644); err != nil {
			t.Fatal(err)
		}

		_, err := prepareTarget(path)
		if err == nil || !strings.Contains(err.Error(), "target directory is not empty") {
			t.Fatalf("prepareTarget() error = %v", err)
		}
	})
}

func TestCreateProjectCopiesTemplate(t *testing.T) {
	templateRoot := t.TempDir()
	if err := os.Mkdir(filepath.Join(templateRoot, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(templateRoot, "src", "main.go"), []byte("package main\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(t.TempDir(), "hello")
	got, err := createProject(templates.Template{Name: "go", Path: templateRoot}, target)
	if err != nil {
		t.Fatalf("createProject() error = %v", err)
	}

	wantPath, err := filepath.Abs(target)
	if err != nil {
		t.Fatal(err)
	}
	wantPath = filepath.Clean(wantPath)

	if got.Name != "hello" {
		t.Fatalf("Project.Name = %q, want %q", got.Name, "hello")
	}
	if got.Path != wantPath {
		t.Fatalf("Project.Path = %q, want %q", got.Path, wantPath)
	}

	data, err := os.ReadFile(filepath.Join(target, "src", "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "package main\n" {
		t.Fatalf("copied file = %q", string(data))
	}
}

func TestCreateProjectRejectsNonEmptyTarget(t *testing.T) {
	templateRoot := t.TempDir()
	target := t.TempDir()
	if err := os.WriteFile(filepath.Join(target, "existing.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := createProject(templates.Template{Name: "go", Path: templateRoot}, target)
	if err == nil || !strings.Contains(err.Error(), "target directory is not empty") {
		t.Fatalf("createProject() error = %v", err)
	}
}
