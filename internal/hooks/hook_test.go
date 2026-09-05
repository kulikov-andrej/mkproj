package hooks

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/data/project"
	"github.com/kulikov-andrej/mkproj/internal/data/templates"
)

func runTestHook(
	t *testing.T,
	source string,
) (
	string,
	*bytes.Buffer,
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

	if err := os.WriteFile(
		filepath.Join(
			metadataPath,
			"hook.star",
		),
		[]byte(source),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		project.Project{
			Name: "hello",
			Path: projectPath,
		},
		templates.Template{
			Name: "example",
		},
		strings.NewReader(""),
		&stdout,
		&stderr,
	)

	return projectPath, &stdout, err
}

func TestRunWithoutHook(t *testing.T) {
	projectPath := t.TempDir()

	err := Run(
		project.Project{
			Name: "hello",
			Path: projectPath,
		},
		templates.Template{
			Name: "example",
		},
		strings.NewReader(""),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRun(t *testing.T) {
	projectPath, stdout, err := runTestHook(
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

func TestReplace(t *testing.T) {
	projectPath, _, err := runTestHook(
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
	projectPath, _, err := runTestHook(
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
func TestRunReturnsHookError(t *testing.T) {
	_, _, err := runTestHook(
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
		"template hook failed",
	) {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}
