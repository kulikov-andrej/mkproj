package scripts_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/scripts"
	"go.starlark.net/starlark"
)

func loadTestScript(
	t *testing.T,
	source string,
) (
	string,
	*bytes.Buffer,
	starlark.StringDict,
	error,
) {
	t.Helper()

	projectPath := t.TempDir()

	metadataPath := filepath.Join(
		projectPath,
		".mkproj",
	)

	if err := os.MkdirAll(
		metadataPath,
		0o755,
	); err != nil {
		t.Fatal(err)
	}

	scriptPath := filepath.Join(
		metadataPath,
		"test.star",
	)

	if err := os.WriteFile(
		scriptPath,
		[]byte(source),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	globals, err := scripts.Load(
		scriptPath,
		scripts.Context{
			ProjectName:  "hello",
			ProjectPath:  projectPath,
			TemplateName: "example",
			Stdin:        strings.NewReader(""),
			Stdout:       &stdout,
			Stderr:       &stderr,
		},
	)

	return projectPath, &stdout, globals, err
}

func TestLoad(t *testing.T) {
	projectPath, stdout, _, err := loadTestScript(
		t,
		`
print(project.name)
print(template.name)

write(
    "generated.txt",
    project.name + ":" + template.name,
)
`,
	)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := stdout.String(), "hello\nexample\n"; got != want {
		t.Fatalf(
			"expected stdout %q, got %q",
			want,
			got,
		)
	}

	content, err := os.ReadFile(
		filepath.Join(
			projectPath,
			"generated.txt",
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(content), "hello:example"; got != want {
		t.Fatalf(
			"expected %q, got %q",
			want,
			got,
		)
	}
}

func TestLoadReturnsGlobals(t *testing.T) {
	_, _, globals, err := loadTestScript(
		t,
		`
def build():
    pass
`,
	)
	if err != nil {
		t.Fatal(err)
	}

	value, ok := globals["build"]
	if !ok {
		t.Fatal(`expected global "build"`)
	}

	if _, ok := value.(starlark.Callable); !ok {
		t.Fatalf(
			"expected build to be callable, got %s",
			value.Type(),
		)
	}
}

func TestReplace(t *testing.T) {
	projectPath, _, _, err := loadTestScript(
		t,
		`
write("project.txt", "name={{NAME}}")
count = replace(
    "project.txt",
    "{{NAME}}",
    project.name,
)

if count != 1:
    fail("unexpected replacement count")
`,
	)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(
		filepath.Join(
			projectPath,
			"project.txt",
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(content), "name=hello"; got != want {
		t.Fatalf(
			"expected %q, got %q",
			want,
			got,
		)
	}
}

func TestWriteCannotEscapeProject(t *testing.T) {
	projectPath, _, _, err := loadTestScript(
		t,
		`
write("../outside.txt", "nope")
`,
	)

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(
		err.Error(),
		"escapes",
	) {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	outside := filepath.Join(
		filepath.Dir(projectPath),
		"outside.txt",
	)

	if _, statErr := os.Stat(outside); !os.IsNotExist(statErr) {
		t.Fatalf(
			"outside file should not exist: %v",
			statErr,
		)
	}
}

func TestLoadReturnsScriptError(t *testing.T) {
	_, _, _, err := loadTestScript(
		t,
		`
missing_function()
`,
	)

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(
		err.Error(),
		"missing_function",
	) {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestCall(t *testing.T) {
	projectPath := t.TempDir()

	metadataPath := filepath.Join(
		projectPath,
		".mkproj",
	)
	if err := os.MkdirAll(metadataPath, 0o755); err != nil {
		t.Fatal(err)
	}

	scriptPath := filepath.Join(
		metadataPath,
		"workflow.star",
	)
	if err := os.WriteFile(
		scriptPath,
		[]byte(`
def hello():
    print("workflow works")
    write("result.txt", "ok")
`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	ctx := scripts.Context{
		ProjectName: "hello",
		ProjectPath: projectPath,
		Stdin:       strings.NewReader(""),
		Stdout:      &stdout,
		Stderr:      &stderr,
	}

	globals, err := scripts.Load(scriptPath, ctx)
	if err != nil {
		t.Fatal(err)
	}

	value, ok := globals["hello"]
	if !ok {
		t.Fatal(`expected global "hello"`)
	}

	fn, ok := value.(starlark.Callable)
	if !ok {
		t.Fatalf("expected hello to be callable, got %s", value.Type())
	}

	if err := scripts.Call(fn, ctx); err != nil {
		t.Fatal(err)
	}

	if got, want := stdout.String(), "workflow works\n"; got != want {
		t.Fatalf("expected stdout %q, got %q", want, got)
	}

	content, err := os.ReadFile(filepath.Join(projectPath, "result.txt"))
	if err != nil {
		t.Fatal(err)
	}

	if got, want := string(content), "ok"; got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
