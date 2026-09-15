package cli

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/buildinfo"
)

func TestRunWithoutArgsShowsHelp(t *testing.T) {
	got := runCLI()

	if got.err != nil {
		t.Fatal(got.err)
	}

	if !strings.Contains(got.stdout, "Usage:") {
		t.Fatalf("expected help output, got %q", got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunHelp(t *testing.T) {
	got := runCLI("help")

	if got.err != nil {
		t.Fatal(got.err)
	}

	if !strings.Contains(got.stdout, "Usage:") {
		t.Fatalf("expected help output, got %q", got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunVersion(t *testing.T) {
	got := runCLI("version")

	if got.err != nil {
		t.Fatal(got.err)
	}

	want := fmt.Sprintf("mkproj %s\n", buildinfo.Version)
	if got.stdout != want {
		t.Fatalf("expected stdout %q, got %q", want, got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunHelpProject(t *testing.T) {
	got := runCLI("help", "project")

	if got.err != nil {
		t.Fatal(got.err)
	}

	if !strings.Contains(
		got.stdout,
		"mkproj [<path>] -t <template> [--open]",
	) {
		t.Fatalf("expected project help, got %q", got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunHelpTemplate(t *testing.T) {
	got := runCLI("help", "template")

	if got.err != nil {
		t.Fatal(got.err)
	}

	if !strings.Contains(
		got.stdout,
		"mkproj template <command>",
	) {
		t.Fatalf("expected template help, got %q", got.stdout)
	}

	if got.stderr != "" {
		t.Fatalf("unexpected stderr: %q", got.stderr)
	}
}

func TestRunHelpUnknownTopic(t *testing.T) {
	got := runCLI("help", "what")

	if got.err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(
		got.err.Error(),
		`unknown help topic "what"`,
	) {
		t.Fatalf("unexpected error: %v", got.err)
	}
}

func TestRunHelpRejectsExtraArguments(t *testing.T) {
	got := runCLI("help", "project", "what")

	if got.err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(
		got.err.Error(),
		`unexpected argument "what"`,
	) {
		t.Fatalf("unexpected error: %v", got.err)
	}
}
