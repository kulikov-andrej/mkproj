package templates

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func useTemplateRoot(t *testing.T, root string) {
	t.Helper()

	old := resolveRoot

	resolveRoot = func() (string, error) {
		return root, nil
	}

	t.Cleanup(func() {
		resolveRoot = old
	})
}

func TestList(t *testing.T) {
	root := t.TempDir()
	useTemplateRoot(t, root)

	for _, name := range []string{
		"cpp",
		"go",
	} {
		if err := os.Mkdir(
			filepath.Join(root, name),
			0o755,
		); err != nil {
			t.Fatal(err)
		}
	}

	if err := os.WriteFile(
		filepath.Join(root, "README.txt"),
		[]byte("ignore me"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	got, err := List()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf(
			"expected 2 templates, got %d",
			len(got),
		)
	}

	if got[0].Name != "cpp" ||
		got[1].Name != "go" {
		t.Fatalf(
			"unexpected templates: %#v",
			got,
		)
	}
}

func TestGet(t *testing.T) {
	root := t.TempDir()
	useTemplateRoot(t, root)

	templatePath := filepath.Join(
		root,
		"example",
	)

	if err := os.Mkdir(
		templatePath,
		0o755,
	); err != nil {
		t.Fatal(err)
	}

	got, err := Get("example")
	if err != nil {
		t.Fatal(err)
	}

	if got.Name != "example" {
		t.Fatalf(
			"expected name %q, got %q",
			"example",
			got.Name,
		)
	}

	if got.Path != templatePath {
		t.Fatalf(
			"expected path %q, got %q",
			templatePath,
			got.Path,
		)
	}
}

func TestGetMissingTemplate(t *testing.T) {
	root := t.TempDir()
	useTemplateRoot(t, root)

	_, err := Get("missing")

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(
		err.Error(),
		"not found",
	) {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}

func TestGetRejectsInvalidName(t *testing.T) {
	root := t.TempDir()
	useTemplateRoot(t, root)

	tests := []string{
		"",
		".",
		"..",
		"../outside",
		"foo/bar",
		`foo\bar`,
	}

	for _, name := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Get(name)

			if err == nil {
				t.Fatalf(
					"expected %q to be rejected",
					name,
				)
			}
		})
	}
}

func TestListReturnsRootError(t *testing.T) {
	old := resolveRoot

	expected := errors.New(
		"root unavailable",
	)

	resolveRoot = func() (string, error) {
		return "", expected
	}

	t.Cleanup(func() {
		resolveRoot = old
	})

	_, err := List()

	if !errors.Is(err, expected) {
		t.Fatalf(
			"expected %v, got %v",
			expected,
			err,
		)
	}
}
