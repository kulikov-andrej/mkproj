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

func TestRunWorkflowHelpOutsideProject(t *testing.T) {
	resetDependencies(t)

	getCurrentProject = func() (project.Project, error) {
		return project.Project{}, errors.New("not a project")
	}

	listWorkflowCommands = func(
		project.Project,
		io.Reader,
		io.Writer,
		io.Writer,
	) ([]string, error) {
		t.Fatal("workflow commands should not be loaded")
		return nil, nil
	}

	got := runCLI("run")
	if got.err != nil {
		t.Fatal(got.err)
	}

	if got.stdout != "Usage:\n  mkproj run <command>\n" {
		t.Fatalf("unexpected stdout: %q", got.stdout)
	}
}

func TestRunWorkflowHelpShowsProjectCommands(t *testing.T) {
	resetDependencies(t)

	proj := project.Project{
		Name: "hello",
		Path: "project/path",
	}

	getCurrentProject = func() (project.Project, error) {
		return proj, nil
	}

	listWorkflowCommands = func(
		gotProj project.Project,
		_ io.Reader,
		_ io.Writer,
		_ io.Writer,
	) ([]string, error) {
		if gotProj != proj {
			t.Fatalf("expected project %#v, got %#v", proj, gotProj)
		}

		return []string{"build", "test"}, nil
	}

	got := runCLI("run")
	if got.err != nil {
		t.Fatal(got.err)
	}

	want := "Usage:\n  mkproj run <command>\n\nCommands:\n  build\n  test\n"
	if got.stdout != want {
		t.Fatalf("expected stdout %q, got %q", want, got.stdout)
	}
}

func TestRunWorkflowCommand(t *testing.T) {
	resetDependencies(t)

	proj := project.Project{
		Name: "hello",
		Path: "project/path",
	}

	getCurrentProject = func() (project.Project, error) {
		return proj, nil
	}

	called := false

	runProjectWorkflow = func(
		gotProj project.Project,
		command string,
		_ io.Reader,
		_ io.Writer,
		_ io.Writer,
	) error {
		called = true

		if gotProj != proj {
			t.Fatalf("expected project %#v, got %#v", proj, gotProj)
		}

		if command != "debug" {
			t.Fatalf("expected command %q, got %q", "debug", command)
		}

		return nil
	}

	got := runCLI("run", "debug")
	if got.err != nil {
		t.Fatal(got.err)
	}

	if !called {
		t.Fatal("project workflow was not called")
	}
}

func TestRunWorkflowCommandRequiresProject(t *testing.T) {
	resetDependencies(t)

	wantErr := errors.New("not a project")

	getCurrentProject = func() (project.Project, error) {
		return project.Project{}, wantErr
	}

	runProjectWorkflow = func(
		project.Project,
		string,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		t.Fatal("project workflow should not be called")
		return nil
	}

	got := runCLI("run", "debug")
	if !errors.Is(got.err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, got.err)
	}
}

func TestRunWorkflowRejectsExtraArguments(t *testing.T) {
	resetDependencies(t)

	getCurrentProject = func() (project.Project, error) {
		t.Fatal("project should not be resolved")
		return project.Project{}, nil
	}

	got := runCLI("run", "debug", "extra")
	if got.err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(got.err.Error(), `unexpected argument "extra"`) {
		t.Fatalf("unexpected error: %v", got.err)
	}
}
