package project

import (
	"bytes"
	"errors"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/scripts"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
	"go.starlark.net/starlark"
)

func resetDependencies(t *testing.T) {
	t.Helper()

	oldGetTemplate := getTemplate
	oldCreateProject := createProject
	oldGetCurrentProject := getCurrentProject
	oldLoadScript := loadScript
	oldCallScript := callScript
	oldCleanupMetadata := cleanupMetadata
	oldOpenEditor := openEditor

	t.Cleanup(func() {
		getTemplate = oldGetTemplate
		createProject = oldCreateProject
		getCurrentProject = oldGetCurrentProject
		loadScript = oldLoadScript
		callScript = oldCallScript
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

	stdin := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

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

	loadScript = func(
		path string,
		ctx scripts.Context,
	) (starlark.StringDict, error) {
		calls = append(calls, "setup")

		wantPath := filepath.Join(
			proj.Path,
			".mkproj",
			"setup.star",
		)

		if path != wantPath {
			t.Fatalf(
				"expected script path %q, got %q",
				wantPath,
				path,
			)
		}

		if ctx.ProjectName != proj.Name {
			t.Fatalf(
				"expected project name %q, got %q",
				proj.Name,
				ctx.ProjectName,
			)
		}

		if ctx.ProjectPath != proj.Path {
			t.Fatalf(
				"expected project path %q, got %q",
				proj.Path,
				ctx.ProjectPath,
			)
		}

		if ctx.TemplateName != tmpl.Name {
			t.Fatalf(
				"expected template name %q, got %q",
				tmpl.Name,
				ctx.TemplateName,
			)
		}

		if ctx.Stdin != stdin {
			t.Fatal("unexpected stdin")
		}

		if ctx.Stdout != stdout {
			t.Fatal("unexpected stdout")
		}

		if ctx.Stderr != stderr {
			t.Fatal("unexpected stderr")
		}

		return starlark.StringDict{}, nil
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
		stdin,
		stdout,
		stderr,
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

	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf(
			"expected calls %v, got %v",
			wantCalls,
			calls,
		)
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

	loadScript = func(
		string,
		scripts.Context,
	) (starlark.StringDict, error) {
		t.Fatal("setup should not be called")
		return nil, nil
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

	loadScript = func(
		string,
		scripts.Context,
	) (starlark.StringDict, error) {
		return nil, setupErr
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

	loadScript = func(
		string,
		scripts.Context,
	) (starlark.StringDict, error) {
		return starlark.StringDict{}, nil
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

func TestWorkflowCommands(t *testing.T) {
	resetDependencies(t)

	proj := Project{
		Name: "hello",
		Path: "target",
	}

	loadScript = func(
		path string,
		ctx scripts.Context,
	) (starlark.StringDict, error) {
		wantPath := filepath.Join(
			proj.Path,
			".mkproj",
			"workflow.star",
		)
		if path != wantPath {
			t.Fatalf("expected path %q, got %q", wantPath, path)
		}

		if ctx.ProjectName != proj.Name || ctx.ProjectPath != proj.Path {
			t.Fatalf("unexpected context: %#v", ctx)
		}

		return starlark.StringDict{
			"test":  starlark.NewBuiltin("test", nil),
			"value": starlark.String("not a command"),
			"build": starlark.NewBuiltin("build", nil),
		}, nil
	}

	commands, err := WorkflowCommands(
		proj,
		strings.NewReader(""),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"build", "test"}
	if !reflect.DeepEqual(commands, want) {
		t.Fatalf("expected commands %v, got %v", want, commands)
	}
}

func TestRunWorkflow(t *testing.T) {
	resetDependencies(t)

	proj := Project{
		Name: "hello",
		Path: "target",
	}

	stdin := strings.NewReader("")
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	fn := starlark.NewBuiltin("build", nil)

	loadScript = func(
		path string,
		ctx scripts.Context,
	) (starlark.StringDict, error) {
		wantPath := filepath.Join(
			proj.Path,
			".mkproj",
			"workflow.star",
		)
		if path != wantPath {
			t.Fatalf("expected path %q, got %q", wantPath, path)
		}

		if ctx.Stdin != stdin || ctx.Stdout != stdout || ctx.Stderr != stderr {
			t.Fatal("unexpected workflow streams")
		}

		return starlark.StringDict{
			"build": fn,
		}, nil
	}

	called := false

	callScript = func(
		gotFn starlark.Callable,
		ctx scripts.Context,
	) error {
		called = true

		if gotFn != fn {
			t.Fatal("unexpected callable")
		}

		if ctx.ProjectName != proj.Name || ctx.ProjectPath != proj.Path {
			t.Fatalf("unexpected context: %#v", ctx)
		}

		return nil
	}

	if err := RunWorkflow(
		proj,
		"build",
		stdin,
		stdout,
		stderr,
	); err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("script callable was not called")
	}
}

func TestRunWorkflowCommandNotFound(t *testing.T) {
	resetDependencies(t)

	loadScript = func(
		string,
		scripts.Context,
	) (starlark.StringDict, error) {
		return starlark.StringDict{}, nil
	}

	err := RunWorkflow(
		Project{Path: "target"},
		"missing",
		strings.NewReader(""),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), `workflow command "missing" not found`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunWorkflowCommandMustBeCallable(t *testing.T) {
	resetDependencies(t)

	loadScript = func(
		string,
		scripts.Context,
	) (starlark.StringDict, error) {
		return starlark.StringDict{
			"build": starlark.String("nope"),
		}, nil
	}

	err := RunWorkflow(
		Project{Path: "target"},
		"build",
		strings.NewReader(""),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)
	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(err.Error(), `workflow command "build" is not callable`) {
		t.Fatalf("unexpected error: %v", err)
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
