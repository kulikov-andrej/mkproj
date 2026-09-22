package project

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
)

func writeWorkflow(t *testing.T, source string) Project {
	t.Helper()

	path := t.TempDir()
	if err := os.Mkdir(filepath.Join(path, ".mkproj"), 0o755); err != nil {
		t.Fatal(err)
	}

	workflowPath := filepath.Join(path, ".mkproj", "workflow.star")
	if err := os.WriteFile(workflowPath, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}

	return Project{
		Name: "hello",
		Path: path,
	}
}

func TestWorkflowCommands(t *testing.T) {
	proj := writeWorkflow(t, `
def test():
    pass

value = "not a command"

def build():
    pass
`)

	commands, err := WorkflowCommands(proj, libio.IO{})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"build", "test"}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("WorkflowCommands() = %v, want %v", commands, want)
	}
}

func TestRunWorkflow(t *testing.T) {
	proj := writeWorkflow(t, `
def build():
    print("workflow works")
    write("result.txt", project.name)
`)
	var stdout bytes.Buffer

	if err := RunWorkflow(proj, "build", libio.IO{Out: &stdout}); err != nil {
		t.Fatal(err)
	}

	if got, want := stdout.String(), "workflow works\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}

	content, err := os.ReadFile(filepath.Join(proj.Path, "result.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(content), "hello"; got != want {
		t.Fatalf("result.txt = %q, want %q", got, want)
	}
}

func TestRunWorkflowCommandNotFound(t *testing.T) {
	proj := writeWorkflow(t, "")

	err := RunWorkflow(proj, "missing", libio.IO{})
	if err == nil || !strings.Contains(err.Error(), `workflow command "missing" not found`) {
		t.Fatalf("RunWorkflow() error = %v", err)
	}
}

func TestRunWorkflowCommandMustBeCallable(t *testing.T) {
	proj := writeWorkflow(t, `build = "nope"`)

	err := RunWorkflow(proj, "build", libio.IO{})
	if err == nil || !strings.Contains(err.Error(), `workflow command "build" is not callable`) {
		t.Fatalf("RunWorkflow() error = %v", err)
	}
}
