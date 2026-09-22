package project

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func TestCreate(t *testing.T) {
	templateRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(templateRoot, "README.md"), []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	old := findTemplate
	findTemplate = func(name string) (templates.Template, error) {
		if name != "go" {
			t.Fatalf("findTemplate() name = %q, want %q", name, "go")
		}
		return templates.Template{Name: name, Path: templateRoot}, nil
	}
	t.Cleanup(func() { findTemplate = old })

	target := filepath.Join(t.TempDir(), "app")
	proj, err := Create("go", target, libio.IO{})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if proj.Name != "app" {
		t.Fatalf("Project.Name = %q, want %q", proj.Name, "app")
	}

	data, err := os.ReadFile(filepath.Join(target, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "hello" {
		t.Fatalf("copied file = %q, want %q", string(data), "hello")
	}
}

func TestCreateReturnsTemplateError(t *testing.T) {
	wantErr := errors.New("template failed")
	old := findTemplate
	findTemplate = func(string) (templates.Template, error) {
		return templates.Template{}, wantErr
	}
	t.Cleanup(func() { findTemplate = old })

	_, err := Create("go", "target", libio.IO{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("Create() error = %v, want %v", err, wantErr)
	}
}
