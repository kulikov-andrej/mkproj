package cli

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/data/project"
	"github.com/kulikov-andrej/mkproj/internal/data/templates"
)

func resetDependencies(t *testing.T) {
	t.Helper()

	oldListTemplates := listTemplates
	oldCreateProject := createProject
	oldOpenProject := openProject

	t.Cleanup(func() {
		listTemplates = oldListTemplates
		createProject = oldCreateProject
		openProject = oldOpenProject
	})
}

func TestRunWithoutArgsShowsHelp(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		nil,
		strings.NewReader(""),
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(
		stdout.String(),
		"Usage:",
	) {
		t.Fatalf(
			"expected help output, got %q",
			stdout.String(),
		)
	}

	if stderr.Len() != 0 {
		t.Fatalf(
			"unexpected stderr: %q",
			stderr.String(),
		)
	}
}
func TestRunHelp(t *testing.T) {
	var stdout bytes.Buffer

	err := Run(
		[]string{"--help"},
		strings.NewReader(""),
		&stdout,
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(
		stdout.String(),
		"Usage:",
	) {
		t.Fatalf(
			"expected help output, got %q",
			stdout.String(),
		)
	}
}
func TestRunList(t *testing.T) {
	resetDependencies(t)

	listTemplates = func() ([]templates.Template, error) {
		return []templates.Template{
			{Name: "cpp"},
			{Name: "go"},
		}, nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{"--list"},
		strings.NewReader(""),
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := stdout.String(), "cpp\ngo\n"; got != want {
		t.Fatalf(
			"expected stdout %q, got %q",
			want,
			got,
		)
	}

	if stderr.Len() != 0 {
		t.Fatalf(
			"unexpected stderr: %q",
			stderr.String(),
		)
	}
}
func TestRunEmptyList(t *testing.T) {
	resetDependencies(t)

	listTemplates = func() ([]templates.Template, error) {
		return nil, nil
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{"--list"},
		strings.NewReader(""),
		&stdout,
		&stderr,
	)
	if err != nil {
		t.Fatal(err)
	}

	if stdout.Len() != 0 {
		t.Fatalf(
			"unexpected stdout: %q",
			stdout.String(),
		)
	}

	if got, want := stderr.String(), "No templates found.\n"; got != want {
		t.Fatalf(
			"expected stderr %q, got %q",
			want,
			got,
		)
	}
}
func TestRunCreate(t *testing.T) {
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

	var stdout bytes.Buffer

	err := Run(
		[]string{
			"--template=example",
			"hello",
		},
		strings.NewReader(""),
		&stdout,
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if gotTemplate != "example" {
		t.Fatalf(
			"expected template %q, got %q",
			"example",
			gotTemplate,
		)
	}

	if gotTarget != "hello" {
		t.Fatalf(
			"expected target %q, got %q",
			"hello",
			gotTarget,
		)
	}

	output := stdout.String()

	if !strings.Contains(
		output,
		`Creating project "hello" using template "example"...`,
	) {
		t.Fatalf(
			"missing create message: %q",
			output,
		)
	}

	if !strings.Contains(
		output,
		`Created "hello".`,
	) {
		t.Fatalf(
			"missing created message: %q",
			output,
		)
	}
}
func TestRunOpen(t *testing.T) {
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

	openProject = func(
		proj project.Project,
	) error {
		opened = proj
		return nil
	}

	var stdout bytes.Buffer

	err := Run(
		[]string{
			"hello",
			"--template=example",
			"--open",
		},
		strings.NewReader(""),
		&stdout,
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if opened != expected {
		t.Fatalf(
			"expected %#v, got %#v",
			expected,
			opened,
		)
	}

	if !strings.Contains(
		stdout.String(),
		"Opening Code...\n",
	) {
		t.Fatalf(
			"missing open message: %q",
			stdout.String(),
		)
	}
}
func TestRunDoesNotOpenAfterCreateError(t *testing.T) {
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

	openProject = func(
		project.Project,
	) error {
		t.Fatal("open should not be called")
		return nil
	}

	err := Run(
		[]string{
			"hello",
			"--template=example",
			"--open",
		},
		strings.NewReader(""),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if !errors.Is(err, expected) {
		t.Fatalf(
			"expected %v, got %v",
			expected,
			err,
		)
	}
}
