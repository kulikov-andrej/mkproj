package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func TestRunWithoutArgsShowsHelp(t *testing.T) {
	streams, out, errOut := testStreams()

	if err := Run(nil, streams); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(out.String(), "Usage:") || !strings.Contains(out.String(), "mkproj <command>") {
		t.Fatalf("stdout = %q, want main help", out.String())
	}
	if errOut.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", errOut.String())
	}
}

func testStreams() (libio.IO, *bytes.Buffer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	errOut := &bytes.Buffer{}
	return libio.IO{
		In:  &bytes.Buffer{},
		Out: out,
		Err: errOut,
	}, out, errOut
}

func TestRunDispatchesHelp(t *testing.T) {
	streams, out, _ := testStreams()

	if err := Run([]string{"help", "project"}, streams); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(out.String(), "Project path") {
		t.Fatalf("stdout = %q, want project help", out.String())
	}
}

func TestRunDispatchesTemplate(t *testing.T) {
	old := listTemplates
	listTemplates = func() ([]templates.Template, error) {
		return []templates.Template{{Name: "go"}}, nil
	}
	t.Cleanup(func() { listTemplates = old })

	streams, out, _ := testStreams()
	if err := Run([]string{"template", "list"}, streams); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if got, want := out.String(), "go\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunDispatchesProject(t *testing.T) {
	old := createProject
	createProject = func(template, target string, streams libio.IO) (project.Project, error) {
		return project.Project{Name: "app", Path: "app"}, nil
	}
	t.Cleanup(func() { createProject = old })

	streams, out, _ := testStreams()
	if err := Run([]string{"app", "-t", "go"}, streams); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !strings.Contains(out.String(), `Created "app".`) {
		t.Fatalf("stdout = %q, want project creation output", out.String())
	}
}
