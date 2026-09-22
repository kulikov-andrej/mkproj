package cli

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project"
)

func TestParseProjectArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want projectArgs
	}{
		{name: "short template", args: []string{"app", "-t", "go"}, want: projectArgs{target: "app", template: "go"}},
		{name: "long template", args: []string{"--template", "go", "app"}, want: projectArgs{target: "app", template: "go"}},
		{name: "template equals", args: []string{"--template=go", "app"}, want: projectArgs{target: "app", template: "go"}},
		{name: "short open", args: []string{"-o", "-t", "go"}, want: projectArgs{template: "go", open: true}},
		{name: "long open", args: []string{"--open", "-t", "go"}, want: projectArgs{template: "go", open: true}},
		{name: "options after target", args: []string{"app", "--open", "--template=go"}, want: projectArgs{target: "app", template: "go", open: true}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseProjectArgs(tt.args)
			if err != nil {
				t.Fatalf("parseProjectArgs() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("parseProjectArgs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestParseProjectArgsErrors(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "missing short template", args: []string{"-t"}, wantErr: "-t requires a value"},
		{name: "missing long template", args: []string{"--template"}, wantErr: "--template requires a value"},
		{name: "blank short template", args: []string{"-t", "   "}, wantErr: "-t requires a value"},
		{name: "blank equals template", args: []string{"--template=   "}, wantErr: "--template requires a value"},
		{name: "unknown option", args: []string{"--wat"}, wantErr: "unknown option: --wat"},
		{name: "second target", args: []string{"one", "two", "-t", "go"}, wantErr: "unexpected argument: two"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseProjectArgs(tt.args)
			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("parseProjectArgs() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestRunProject(t *testing.T) {
	var gotTemplate, gotTarget string
	oldCreate := createProject
	oldOpen := openProject
	createProject = func(template, target string, streams libio.IO) (project.Project, error) {
		gotTemplate = template
		gotTarget = target
		return project.Project{Name: "app", Path: "C:/work/app"}, nil
	}
	openCalled := false
	openProject = func(project.Project) error {
		openCalled = true
		return nil
	}
	t.Cleanup(func() {
		createProject = oldCreate
		openProject = oldOpen
	})

	streams, out, _ := testStreams()
	if err := runProject([]string{"app", "-t", "go"}, streams); err != nil {
		t.Fatalf("runProject() error = %v", err)
	}

	if gotTemplate != "go" || gotTarget != "app" {
		t.Fatalf("createProject() got template=%q target=%q", gotTemplate, gotTarget)
	}
	if openCalled {
		t.Fatal("openProject() called without --open")
	}

	for _, want := range []string{
		`Creating project "app" using template "go"...`,
		`Created "app".`,
		"C:/work/app",
		"Have a vibey coding :)",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, want substring %q", out.String(), want)
		}
	}
}

func TestRunProjectDefaultsTargetToCurrentDirectory(t *testing.T) {
	old := createProject
	createProject = func(template, target string, streams libio.IO) (project.Project, error) {
		if target != "." {
			t.Fatalf("target = %q, want %q", target, ".")
		}
		return project.Project{Name: "app", Path: "."}, nil
	}
	t.Cleanup(func() { createProject = old })

	streams, _, _ := testStreams()
	if err := runProject([]string{"-t", "go"}, streams); err != nil {
		t.Fatalf("runProject() error = %v", err)
	}
}

func TestRunProjectRequiresTemplate(t *testing.T) {
	streams, _, _ := testStreams()

	err := runProject([]string{"app"}, streams)
	if err == nil || err.Error() != "template is required" {
		t.Fatalf("runProject() error = %v", err)
	}
}

func TestRunProjectReturnsCreateError(t *testing.T) {
	wantErr := errors.New("create failed")
	old := createProject
	createProject = func(string, string, libio.IO) (project.Project, error) {
		return project.Project{}, wantErr
	}
	t.Cleanup(func() { createProject = old })

	streams, _, _ := testStreams()
	err := runProject([]string{"-t", "go"}, streams)
	if !errors.Is(err, wantErr) {
		t.Fatalf("runProject() error = %v, want %v", err, wantErr)
	}
}

func TestRunProjectOpen(t *testing.T) {
	oldCreate := createProject
	oldOpen := openProject
	createProject = func(string, string, libio.IO) (project.Project, error) {
		return project.Project{Name: "app", Path: "app"}, nil
	}
	gotProject := project.Project{}
	openProject = func(proj project.Project) error {
		gotProject = proj
		return nil
	}
	t.Cleanup(func() {
		createProject = oldCreate
		openProject = oldOpen
	})

	streams, out, _ := testStreams()
	if err := runProject([]string{"-t", "go", "--open"}, streams); err != nil {
		t.Fatalf("runProject() error = %v", err)
	}
	if gotProject.Name != "app" || gotProject.Path != "app" {
		t.Fatalf("openProject() project = %#v", gotProject)
	}
	if !strings.Contains(out.String(), "Opening Code...") {
		t.Fatalf("stdout = %q, want opening message", out.String())
	}
}

func TestRunProjectReturnsOpenError(t *testing.T) {
	wantErr := errors.New("open failed")
	oldCreate := createProject
	oldOpen := openProject
	createProject = func(string, string, libio.IO) (project.Project, error) {
		return project.Project{Name: "app", Path: "app"}, nil
	}
	openProject = func(project.Project) error { return wantErr }
	t.Cleanup(func() {
		createProject = oldCreate
		openProject = oldOpen
	})

	streams, out, _ := testStreams()
	err := runProject([]string{"-t", "go", "--open"}, streams)
	if !errors.Is(err, wantErr) {
		t.Fatalf("runProject() error = %v, want %v", err, wantErr)
	}
	if strings.Contains(out.String(), "Have a vibey coding :)") {
		t.Fatalf("success message printed after open failure: %q", out.String())
	}
}
