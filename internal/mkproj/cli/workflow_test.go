package cli

import (
	"errors"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project"
)

func TestRunWorkflowHelpOutsideProject(t *testing.T) {
	oldGetProject := getProject
	oldListWorkflow := listWorkflow
	getProject = func() (project.Project, error) {
		return project.Project{}, errors.New("not a project")
	}
	listWorkflow = func(project.Project, libio.IO) ([]string, error) {
		t.Fatal("workflow commands should not be loaded")
		return nil, nil
	}
	t.Cleanup(func() {
		getProject = oldGetProject
		listWorkflow = oldListWorkflow
	})

	streams, out, _ := testStreams()
	if err := runWorkflow(nil, streams); err != nil {
		t.Fatal(err)
	}

	want := "Usage:\n  mkproj run <command>\n"
	if got := out.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunWorkflowHelpShowsProjectCommands(t *testing.T) {
	oldGetProject := getProject
	oldListWorkflow := listWorkflow
	proj := project.Project{Name: "hello", Path: "project/path"}
	getProject = func() (project.Project, error) { return proj, nil }
	listWorkflow = func(got project.Project, _ libio.IO) ([]string, error) {
		if got != proj {
			t.Fatalf("project = %#v, want %#v", got, proj)
		}
		return []string{"build", "test"}, nil
	}
	t.Cleanup(func() {
		getProject = oldGetProject
		listWorkflow = oldListWorkflow
	})

	streams, out, _ := testStreams()
	if err := runWorkflow(nil, streams); err != nil {
		t.Fatal(err)
	}

	want := "Usage:\n  mkproj run <command>\n\nCommands:\n  build\n  test\n"
	if got := out.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunWorkflowCommand(t *testing.T) {
	oldGetProject := getProject
	oldRunWorkflow := runWorkflowFn
	proj := project.Project{Name: "hello", Path: "project/path"}
	getProject = func() (project.Project, error) { return proj, nil }
	called := false
	runWorkflowFn = func(got project.Project, command string, _ libio.IO) error {
		called = true
		if got != proj || command != "debug" {
			t.Fatalf("got project=%#v command=%q", got, command)
		}
		return nil
	}
	t.Cleanup(func() {
		getProject = oldGetProject
		runWorkflowFn = oldRunWorkflow
	})

	if err := runWorkflow([]string{"debug"}, libio.IO{}); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("project workflow was not called")
	}
}

func TestRunWorkflowCommandRequiresProject(t *testing.T) {
	oldGetProject := getProject
	oldRunWorkflow := runWorkflowFn
	wantErr := errors.New("not a project")
	getProject = func() (project.Project, error) { return project.Project{}, wantErr }
	runWorkflowFn = func(project.Project, string, libio.IO) error {
		t.Fatal("project workflow should not be called")
		return nil
	}
	t.Cleanup(func() {
		getProject = oldGetProject
		runWorkflowFn = oldRunWorkflow
	})

	err := runWorkflow([]string{"debug"}, libio.IO{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
}

func TestRunWorkflowRejectsExtraArguments(t *testing.T) {
	streams, _, _ := testStreams()
	err := runWorkflow([]string{"debug", "extra"}, streams)
	if err == nil || !strings.Contains(err.Error(), `unexpected argument "extra"`) {
		t.Fatalf("runWorkflow() error = %v", err)
	}
}
