package cli

import (
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func TestRunTemplateShowsHelp(t *testing.T) {
	got := runCLI("template")

	if got.err != nil {
		t.Fatal(got.err)
	}

	if !strings.Contains(got.stdout, "mkproj template <command>") {
		t.Fatalf("expected template help, got %q", got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunTemplateList(t *testing.T) {
	resetDependencies(t)

	listTemplates = func() ([]templates.Template, error) {
		return []templates.Template{
			{Name: "cpp"},
			{Name: "go"},
		}, nil
	}

	got := runCLI("template", "list")

	if got.err != nil {
		t.Fatal(got.err)
	}

	if want := "cpp\ngo\n"; got.stdout != want {
		t.Fatalf("expected stdout %q, got %q", want, got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunTemplateListEmpty(t *testing.T) {
	resetDependencies(t)

	listTemplates = func() ([]templates.Template, error) {
		return nil, nil
	}

	got := runCLI("template", "list")

	if got.err != nil {
		t.Fatal(got.err)
	}

	if got.stdout != "" {
		t.Fatalf("unexpected stdout: %q", got.stdout)
	}

	if want := "No templates found.\n"; got.stderr != want {
		t.Fatalf("expected stderr %q, got %q", want, got.stderr)
	}
}

func TestRunTemplateRejectsUnknownCommand(t *testing.T) {
	got := runCLI("template", "what")

	if got.err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(got.err.Error(), `unknown template command "what"`) {
		t.Fatalf("unexpected error: %v", got.err)
	}
}

func TestRunTemplateListRejectsExtraArguments(t *testing.T) {
	got := runCLI("template", "list", "what")

	if got.err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(got.err.Error(), `unexpected argument "what"`) {
		t.Fatalf("unexpected error: %v", got.err)
	}
}
