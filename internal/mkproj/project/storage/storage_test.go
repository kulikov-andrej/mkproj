package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/model"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
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

	tmpl := templates.Template{
		Name: "example",
		Path: templatePath,
	}

	got, err := Create(
		tmpl,
		target,
	)
	if err != nil {
		t.Fatal(err)
	}

	wantPath, err := filepath.Abs(target)
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

	if got.Path != wantPath {
		t.Fatalf(
			"expected project path %q, got %q",
			wantPath,
			got.Path,
		)
	}

	content, err := os.ReadFile(
		filepath.Join(target, "main.txt"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(content), "hello"; got != want {
		t.Fatalf(
			"expected %q, got %q",
			want,
			got,
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

	tmpl := templates.Template{
		Name: "example",
		Path: templatePath,
	}

	_, err := Create(
		tmpl,
		target,
	)

	if err == nil {
		t.Fatal("expected an error")
	}

	content, readErr := os.ReadFile(existing)
	if readErr != nil {
		t.Fatal(readErr)
	}

	if got, want := string(content), "keep me"; got != want {
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

	tmpl := templates.Template{
		Name: "example",
		Path: templatePath,
	}

	_, err := Create(
		tmpl,
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

	setupPath := filepath.Join(metadata, "setup.star")
	if err := os.WriteFile(
		setupPath,
		[]byte("setup"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	workflowPath := filepath.Join(metadata, "workflow.star")
	if err := os.WriteFile(
		workflowPath,
		[]byte("workflow"),
		0o644,
	); err != nil {
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

	proj := model.Project{
		Name: "example",
		Path: root,
	}

	if err := CleanupMetadata(proj); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(setupPath); !os.IsNotExist(err) {
		t.Fatalf(
			"expected setup metadata to be removed, got %v",
			err,
		)
	}

	if _, err := os.Stat(workflowPath); err != nil {
		t.Fatalf(
			"workflow metadata should remain: %v",
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

func TestGetCurrentProjectFromNestedDirectory(t *testing.T) {
	root := t.TempDir()
	metadata := filepath.Join(root, ".mkproj")
	nested := filepath.Join(root, "src", "nested")

	if err := os.MkdirAll(metadata, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}

	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(nested); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldwd); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	got, err := GetCurrentProject()
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != filepath.Base(root) {
		t.Fatalf("expected name %q, got %q", filepath.Base(root), got.Name)
	}

	if got.Path != root {
		t.Fatalf("expected path %q, got %q", root, got.Path)
	}
}

func TestCleanupMetadataRemovesEmptyMetadataDirectory(t *testing.T) {
	root := t.TempDir()
	metadata := filepath.Join(root, ".mkproj")

	if err := os.MkdirAll(metadata, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(
		filepath.Join(metadata, "setup.star"),
		[]byte("setup"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	if err := CleanupMetadata(model.Project{Path: root}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(metadata); !os.IsNotExist(err) {
		t.Fatalf("expected empty metadata directory to be removed, got %v", err)
	}
}
