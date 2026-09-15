package project

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func resetDependencies(t *testing.T) {
	t.Helper()

	oldGetTemplate := getTemplate
	oldCreateProject := createProject
	oldRunSetup := runSetup
	oldCleanupMetadata := cleanupMetadata
	oldOpenEditor := openEditor

	t.Cleanup(func() {
		getTemplate = oldGetTemplate
		createProject = oldCreateProject
		runSetup = oldRunSetup
		cleanupMetadata = oldCleanupMetadata
		openEditor = oldOpenEditor
	})
}

func TestCreate(t *testing.T) {
	resetDependencies(t)

	tmpl := templates.Template{
		Name: "example",
		Path: "template",
	}

	proj := Project{
		Name: "hello",
		Path: "target",
	}

	var calls []string

	getTemplate = func(name string) (templates.Template, error) {
		calls = append(calls, "get")

		if name != tmpl.Name {
			t.Fatalf(
				"expected template name %q, got %q",
				tmpl.Name,
				name,
			)
		}

		return tmpl, nil
	}

	createProject = func(
		gotTmpl templates.Template,
		target string,
	) (Project, error) {
		calls = append(calls, "create")

		if gotTmpl != tmpl {
			t.Fatalf(
				"expected template %#v, got %#v",
				tmpl,
				gotTmpl,
			)
		}

		if target != "hello" {
			t.Fatalf(
				"expected target %q, got %q",
				"hello",
				target,
			)
		}

		return proj, nil
	}

	runSetup = func(
		gotProj Project,
		gotTmpl templates.Template,
		_ io.Reader,
		_ io.Writer,
		_ io.Writer,
	) error {
		calls = append(calls, "setup")

		if gotProj != proj {
			t.Fatalf(
				"expected project %#v, got %#v",
				proj,
				gotProj,
			)
		}

		if gotTmpl != tmpl {
			t.Fatalf(
				"expected template %#v, got %#v",
				tmpl,
				gotTmpl,
			)
		}

		return nil
	}

	cleanupMetadata = func(gotProj Project) error {
		calls = append(calls, "cleanup")

		if gotProj != proj {
			t.Fatalf(
				"expected project %#v, got %#v",
				proj,
				gotProj,
			)
		}

		return nil
	}

	got, err := Create(
		tmpl.Name,
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

	wantCalls := []string{
		"get",
		"create",
		"setup",
		"cleanup",
	}

	if len(calls) != len(wantCalls) {
		t.Fatalf(
			"expected calls %v, got %v",
			wantCalls,
			calls,
		)
	}

	for i := range wantCalls {
		if calls[i] != wantCalls[i] {
			t.Fatalf(
				"expected calls %v, got %v",
				wantCalls,
				calls,
			)
		}
	}
}

func TestCreateStopsOnTemplateError(t *testing.T) {
	resetDependencies(t)

	wantErr := errors.New("template error")

	getTemplate = func(string) (templates.Template, error) {
		return templates.Template{}, wantErr
	}

	createProject = func(
		templates.Template,
		string,
	) (Project, error) {
		t.Fatal("create should not be called")
		return Project{}, nil
	}

	_, err := Create(
		"example",
		"hello",
		&bytes.Buffer{},
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"expected %v, got %v",
			wantErr,
			err,
		)
	}
}

func TestCreateStopsOnCreateError(t *testing.T) {
	resetDependencies(t)

	wantErr := errors.New("create error")

	tmpl := templates.Template{
		Name: "example",
	}

	getTemplate = func(string) (templates.Template, error) {
		return tmpl, nil
	}

	createProject = func(
		templates.Template,
		string,
	) (Project, error) {
		return Project{}, wantErr
	}

	runSetup = func(
		Project,
		templates.Template,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		t.Fatal("setup should not be called")
		return nil
	}

	_, err := Create(
		tmpl.Name,
		"hello",
		&bytes.Buffer{},
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"expected %v, got %v",
			wantErr,
			err,
		)
	}
}

func TestCreateCleansUpAfterSetupError(t *testing.T) {
	resetDependencies(t)

	setupErr := errors.New("setup error")

	tmpl := templates.Template{
		Name: "example",
	}

	proj := Project{
		Name: "hello",
		Path: "target",
	}

	getTemplate = func(string) (templates.Template, error) {
		return tmpl, nil
	}

	createProject = func(
		templates.Template,
		string,
	) (Project, error) {
		return proj, nil
	}

	runSetup = func(
		Project,
		templates.Template,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		return setupErr
	}

	cleaned := false

	cleanupMetadata = func(Project) error {
		cleaned = true
		return nil
	}

	_, err := Create(
		tmpl.Name,
		"hello",
		&bytes.Buffer{},
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if !errors.Is(err, setupErr) {
		t.Fatalf(
			"expected %v, got %v",
			setupErr,
			err,
		)
	}

	if !cleaned {
		t.Fatal("metadata cleanup was not called")
	}
}

func TestCreateReturnsCleanupError(t *testing.T) {
	resetDependencies(t)

	wantErr := errors.New("cleanup error")

	tmpl := templates.Template{
		Name: "example",
	}

	proj := Project{
		Name: "hello",
		Path: "target",
	}

	getTemplate = func(string) (templates.Template, error) {
		return tmpl, nil
	}

	createProject = func(
		templates.Template,
		string,
	) (Project, error) {
		return proj, nil
	}

	runSetup = func(
		Project,
		templates.Template,
		io.Reader,
		io.Writer,
		io.Writer,
	) error {
		return nil
	}

	cleanupMetadata = func(Project) error {
		return wantErr
	}

	_, err := Create(
		tmpl.Name,
		"hello",
		&bytes.Buffer{},
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"expected %v, got %v",
			wantErr,
			err,
		)
	}
}

func TestOpen(t *testing.T) {
	resetDependencies(t)

	proj := Project{
		Path: "some/path",
	}

	var gotPath string

	openEditor = func(path string) error {
		gotPath = path
		return nil
	}

	if err := Open(proj); err != nil {
		t.Fatal(err)
	}

	if gotPath != proj.Path {
		t.Fatalf(
			"expected %q, got %q",
			proj.Path,
			gotPath,
		)
	}
}
