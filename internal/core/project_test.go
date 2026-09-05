package core

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/data/project"
	"github.com/kulikov-andrej/mkproj/internal/data/templates"
)

func resetProjectFunctions(t *testing.T) {
	t.Helper()

	oldGetTemplate := getTemplate
	oldCreateProject := createProject
	oldRunHook := runHook
	oldCleanupMetadata := cleanupMetadata
	oldOpenProject := openProject

	t.Cleanup(func() {
		getTemplate = oldGetTemplate
		createProject = oldCreateProject
		runHook = oldRunHook
		cleanupMetadata = oldCleanupMetadata
		openProject = oldOpenProject
	})
}
func TestCreateProject(t *testing.T) {
	resetProjectFunctions(t)

	template := templates.Template{
		Name: "example",
		Path: "template",
	}

	proj := project.Project{
		Name: "hello",
		Path: "target",
	}

	var calls []string

	getTemplate = func(name string) (templates.Template, error) {
		calls = append(calls, "get")

		if name != "example" {
			t.Fatalf("unexpected template name: %q", name)
		}

		return template, nil
	}

	createProject = func(
		gotTemplate templates.Template,
		target string,
	) (project.Project, error) {
		calls = append(calls, "create")

		if gotTemplate != template {
			t.Fatalf("unexpected template: %#v", gotTemplate)
		}

		if target != "hello" {
			t.Fatalf("unexpected target: %q", target)
		}

		return proj, nil
	}

	runHook = func(
		gotProject project.Project,
		gotTemplate templates.Template,
		_ io.Reader,
		_ io.Writer,
		_ io.Writer,
	) error {
		calls = append(calls, "hook")

		if gotProject != proj {
			t.Fatalf("unexpected project: %#v", gotProject)
		}

		if gotTemplate != template {
			t.Fatalf("unexpected template: %#v", gotTemplate)
		}

		return nil
	}

	cleanupMetadata = func(
		gotProject project.Project,
	) error {
		calls = append(calls, "cleanup")

		if gotProject != proj {
			t.Fatalf("unexpected project: %#v", gotProject)
		}

		return nil
	}

	got, err := CreateProject(
		"example",
		"hello",
		&bytes.Buffer{},
		&bytes.Buffer{},
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatal(err)
	}

	if got != proj {
		t.Fatalf(
			"expected %#v, got %#v",
			proj,
			got,
		)
	}

	want := []string{
		"get",
		"create",
		"hook",
		"cleanup",
	}

	if len(calls) != len(want) {
		t.Fatalf(
			"expected calls %v, got %v",
			want,
			calls,
		)
	}

	for i := range want {
		if calls[i] != want[i] {
			t.Fatalf(
				"expected calls %v, got %v",
				want,
				calls,
			)
		}
	}
}
func TestCreateProjectStopsOnTemplateError(t *testing.T) {
	resetProjectFunctions(t)

	expected := errors.New("template error")

	getTemplate = func(
		string,
	) (templates.Template, error) {
		return templates.Template{}, expected
	}

	createProject = func(
		templates.Template,
		string,
	) (project.Project, error) {
		t.Fatal("create should not be called")
		return project.Project{}, nil
	}

	_, err := CreateProject(
		"example",
		"hello",
		&bytes.Buffer{},
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
func TestCreateProjectStopsOnCreateError(t *testing.T) {
	resetProjectFunctions(t)

	expected := errors.New("create error")

	getTemplate = func(
		string,
	) (templates.Template, error) {
		return templates.Template{
			Name: "example",
		}, nil
	}

	createProject = func(
		templates.Template,
		string,
	) (project.Project, error) {
		return project.Project{}, expected
	}

	runHook = func(
		project.Project,
		templates.Template,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		t.Fatal("hook should not be called")
		return nil
	}

	_, err := CreateProject(
		"example",
		"hello",
		&bytes.Buffer{},
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
func TestCreateProjectCleansUpAfterHookError(t *testing.T) {
	resetProjectFunctions(t)

	hookError := errors.New("hook error")

	proj := project.Project{
		Name: "hello",
		Path: "target",
	}

	getTemplate = func(
		string,
	) (templates.Template, error) {
		return templates.Template{
			Name: "example",
		}, nil
	}

	createProject = func(
		templates.Template,
		string,
	) (project.Project, error) {
		return proj, nil
	}

	runHook = func(
		project.Project,
		templates.Template,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		return hookError
	}

	cleaned := false

	cleanupMetadata = func(
		got project.Project,
	) error {
		cleaned = true
		return nil
	}

	_, err := CreateProject(
		"example",
		"hello",
		&bytes.Buffer{},
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if !errors.Is(err, hookError) {
		t.Fatalf(
			"expected %v, got %v",
			hookError,
			err,
		)
	}

	if !cleaned {
		t.Fatal("metadata cleanup was not called")
	}
}
func TestCreateProjectReturnsCleanupError(t *testing.T) {
	resetProjectFunctions(t)

	expected := errors.New("cleanup error")

	proj := project.Project{
		Name: "hello",
		Path: "target",
	}

	getTemplate = func(
		string,
	) (templates.Template, error) {
		return templates.Template{
			Name: "example",
		}, nil
	}

	createProject = func(
		templates.Template,
		string,
	) (project.Project, error) {
		return proj, nil
	}

	runHook = func(
		project.Project,
		templates.Template,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		return nil
	}

	cleanupMetadata = func(
		project.Project,
	) error {
		return expected
	}

	_, err := CreateProject(
		"example",
		"hello",
		&bytes.Buffer{},
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
func TestOpenProject(t *testing.T) {
	resetProjectFunctions(t)

	var got string

	openProject = func(path string) error {
		got = path
		return nil
	}

	proj := project.Project{
		Path: "some/path",
	}

	if err := OpenProject(proj); err != nil {
		t.Fatal(err)
	}

	if got != proj.Path {
		t.Fatalf(
			"expected %q, got %q",
			proj.Path,
			got,
		)
	}
}
