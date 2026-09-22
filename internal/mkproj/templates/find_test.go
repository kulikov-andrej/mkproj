package templates

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestFindTemplate(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "go")
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	stubResolveRoot(t, root, nil)

	got, err := FindTemplate("go")
	if err != nil {
		t.Fatalf("FindTemplate() error = %v", err)
	}

	want := Template{Name: "go", Path: path}
	if got != want {
		t.Fatalf("FindTemplate() = %#v, want %#v", got, want)
	}
}

func TestFindTemplateNotFound(t *testing.T) {
	root := t.TempDir()
	stubResolveRoot(t, root, nil)

	_, err := FindTemplate("missing")
	if err == nil || err.Error() != `template "missing" not found` {
		t.Fatalf("FindTemplate() error = %v", err)
	}
}

func TestFindTemplateRejectsFile(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go"), []byte("file"), 0o644); err != nil {
		t.Fatal(err)
	}
	stubResolveRoot(t, root, nil)

	_, err := FindTemplate("go")
	if err == nil || err.Error() != `template "go" is not a directory` {
		t.Fatalf("FindTemplate() error = %v", err)
	}
}

func TestFindTemplateResolveRootError(t *testing.T) {
	wantErr := errors.New("root failed")
	stubResolveRoot(t, "", wantErr)

	_, err := FindTemplate("go")
	if !errors.Is(err, wantErr) {
		t.Fatalf("FindTemplate() error = %v, want %v", err, wantErr)
	}
}

func TestValidateTemplateName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{name: "valid", input: "go"},
		{name: "hyphenated", input: "go-api"},
		{name: "empty", input: "", wantErr: "template name cannot be empty"},
		{name: "dot", input: ".", wantErr: `invalid template name: "."`},
		{name: "dot dot", input: "..", wantErr: `invalid template name: ".."`},
		{name: "slash", input: "dir/name", wantErr: `invalid template name: "dir/name"`},
		{name: "backslash", input: `dir\name`, wantErr: `invalid template name: "dir\\name"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTemplateName(tt.input)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateTemplateName(%q) error = %v", tt.input, err)
				}
				return
			}

			if err == nil || err.Error() != tt.wantErr {
				t.Fatalf("validateTemplateName(%q) error = %v, want %q", tt.input, err, tt.wantErr)
			}
		})
	}
}

func stubResolveRoot(t *testing.T, root string, err error) {
	t.Helper()
	old := resolveRoot
	resolveRoot = func() (string, error) { return root, err }
	t.Cleanup(func() { resolveRoot = old })
}
