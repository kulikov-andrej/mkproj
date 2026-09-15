package cli

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/project"
)

func TestRunProjectCreate(t *testing.T) {
	resetDependencies(t)

	var gotTemplate string
	var gotTarget string

	createProject = func(
		templateName string,
		target string,
		_ io.Reader,
		_ io.Writer,
		_ io.Writer,
	) (project.Project, error) {
		gotTemplate = templateName
		gotTarget = target

		return project.Project{
			Name: "hello",
			Path: "target/path",
		}, nil
	}

	got := runCLI(
		"--template=example",
		"hello",
	)

	if got.err != nil {
		t.Fatal(got.err)
	}

	if gotTemplate != "example" {
		t.Fatalf("expected template %q, got %q", "example", gotTemplate)
	}

	if gotTarget != "hello" {
		t.Fatalf("expected target %q, got %q", "hello", gotTarget)
	}

	if !strings.Contains(
		got.stdout,
		`Creating project "hello" using template "example"...`,
	) {
		t.Fatalf("missing create message: %q", got.stdout)
	}

	if !strings.Contains(got.stdout, `Created "hello".`) {
		t.Fatalf("missing created message: %q", got.stdout)
	}
}

func TestRunProjectUsesCurrentDirectoryAsDefaultTarget(t *testing.T) {
	resetDependencies(t)

	var gotTarget string

	createProject = func(
		_ string,
		target string,
		_ io.Reader,
		_ io.Writer,
		_ io.Writer,
	) (project.Project, error) {
		gotTarget = target

		return project.Project{
			Name: "project",
			Path: "target/path",
		}, nil
	}

	got := runCLI("-t", "example")
	if got.err != nil {
		t.Fatal(got.err)
	}

	if gotTarget != "." {
		t.Fatalf("expected target %q, got %q", ".", gotTarget)
	}
}

func TestRunProjectOpen(t *testing.T) {
	resetDependencies(t)

	expected := project.Project{
		Name: "hello",
		Path: "target/path",
	}

	createProject = func(
		string,
		string,
		io.Reader,
		io.Writer,
		io.Writer,
	) (project.Project, error) {
		return expected, nil
	}

	var opened project.Project

	openProject = func(proj project.Project) error {
		opened = proj
		return nil
	}

	got := runCLI(
		"hello",
		"--template=example",
		"--open",
	)

	if got.err != nil {
		t.Fatal(got.err)
	}

	if opened != expected {
		t.Fatalf("expected %#v, got %#v", expected, opened)
	}

	if !strings.Contains(got.stdout, "Opening Code...\n") {
		t.Fatalf("missing open message: %q", got.stdout)
	}
}

func TestRunProjectDoesNotOpenAfterCreateError(t *testing.T) {
	resetDependencies(t)

	expected := errors.New("create failed")

	createProject = func(
		string,
		string,
		io.Reader,
		io.Writer,
		io.Writer,
	) (project.Project, error) {
		return project.Project{}, expected
	}

	openProject = func(project.Project) error {
		t.Fatal("open should not be called")
		return nil
	}

	got := runCLI(
		"hello",
		"--template=example",
		"--open",
	)

	if !errors.Is(got.err, expected) {
		t.Fatalf("expected %v, got %v", expected, got.err)
	}
}
