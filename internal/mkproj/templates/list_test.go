package templates

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestListTemplates(t *testing.T) {
	root := t.TempDir()

	for _, name := range []string{"zeta", "alpha", "middle"} {
		if err := os.Mkdir(filepath.Join(root, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(root, "README.txt"), []byte("ignored"), 0o644); err != nil {
		t.Fatal(err)
	}

	stubResolveRoot(t, root, nil)

	got, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}

	want := []Template{
		{Name: "alpha", Path: filepath.Join(root, "alpha")},
		{Name: "middle", Path: filepath.Join(root, "middle")},
		{Name: "zeta", Path: filepath.Join(root, "zeta")},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ListTemplates() = %#v, want %#v", got, want)
	}
}

func TestListTemplatesMissingRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing")
	stubResolveRoot(t, root, nil)

	got, err := ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() error = %v", err)
	}
	if got != nil {
		t.Fatalf("ListTemplates() = %#v, want nil", got)
	}
}

func TestListTemplatesResolveRootError(t *testing.T) {
	wantErr := errors.New("root failed")
	stubResolveRoot(t, "", wantErr)

	_, err := ListTemplates()
	if !errors.Is(err, wantErr) {
		t.Fatalf("ListTemplates() error = %v, want %v", err, wantErr)
	}
}

func TestListTemplatesReadError(t *testing.T) {
	wantErr := errors.New("read failed")

	old := readDir
	readDir = func(string) ([]os.DirEntry, error) {
		return nil, wantErr
	}
	t.Cleanup(func() {
		readDir = old
	})

	stubResolveRoot(t, t.TempDir(), nil)

	_, err := ListTemplates()
	if !errors.Is(err, wantErr) {
		t.Fatalf("ListTemplates() error = %v, want %v", err, wantErr)
	}

	if !strings.Contains(err.Error(), "read templates directory") {
		t.Fatalf("ListTemplates() error = %v, want wrapped read error", err)
	}
}
