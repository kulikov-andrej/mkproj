package project

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func TestRunTemplateSetup(t *testing.T) {
	projectPath := t.TempDir()
	metadataPath := filepath.Join(projectPath, ".mkproj")
	if err := os.Mkdir(metadataPath, 0o755); err != nil {
		t.Fatal(err)
	}

	setupPath := filepath.Join(metadataPath, "setup.star")
	setup := `
print(project.name)
print(template.name)
write("generated.txt", project.name + ":" + template.name)
`
	if err := os.WriteFile(setupPath, []byte(setup), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	err := runTemplateSetup(
		Project{Name: "hello", Path: projectPath},
		templates.Template{Name: "example"},
		libio.IO{Out: &stdout},
	)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := stdout.String(), "hello\nexample\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	content, err := os.ReadFile(filepath.Join(projectPath, "generated.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(content), "hello:example"; got != want {
		t.Fatalf("generated.txt = %q, want %q", got, want)
	}
}

func TestRunTemplateSetupAllowsMissingScript(t *testing.T) {
	err := runTemplateSetup(
		Project{Path: t.TempDir()},
		templates.Template{},
		libio.IO{},
	)
	if err != nil {
		t.Fatalf("runTemplateSetup() error = %v", err)
	}
}

func TestCreateCleansUpSetupAfterError(t *testing.T) {
	templateRoot := t.TempDir()
	metadataPath := filepath.Join(templateRoot, ".mkproj")
	if err := os.Mkdir(metadataPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(metadataPath, "setup.star"),
		[]byte("missing_function()\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	oldFindTemplate := findTemplate
	findTemplate = func(string) (templates.Template, error) {
		return templates.Template{Name: "example", Path: templateRoot}, nil
	}
	t.Cleanup(func() { findTemplate = oldFindTemplate })

	target := filepath.Join(t.TempDir(), "project")
	_, err := Create("example", target, libio.IO{})
	if err == nil || !strings.Contains(err.Error(), "missing_function") {
		t.Fatalf("Create() error = %v", err)
	}

	if _, statErr := os.Stat(filepath.Join(target, ".mkproj", "setup.star")); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("setup.star should be removed, stat error = %v", statErr)
	}
}
