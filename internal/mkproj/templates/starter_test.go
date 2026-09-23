package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallStarter(t *testing.T) {
	root := filepath.Join(t.TempDir(), "templates")
	stubResolveRoot(t, root, nil)

	got, created, err := InstallStarter()
	if err != nil {
		t.Fatalf("InstallStarter() error = %v", err)
	}
	if !created {
		t.Fatal("InstallStarter() created = false, want true")
	}

	wantPath := filepath.Join(root, starterTemplateName)
	if got != (Template{Name: starterTemplateName, Path: wantPath}) {
		t.Fatalf("InstallStarter() = %#v, want starter at %q", got, wantPath)
	}

	files := map[string]string{
		"README.md":                            "{{PROJECT_NAME}}",
		filepath.Join("src", "main.txt"):       "{{PROJECT_NAME}}",
		filepath.Join(".mkproj", "setup.star"): `replace("README.md"`,
	}
	for name, wantText := range files {
		content, err := os.ReadFile(filepath.Join(wantPath, name))
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if !strings.Contains(string(content), wantText) {
			t.Fatalf("%s = %q, want it to contain %q", name, content, wantText)
		}
	}
}

func TestInstallStarterExistingTemplate(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, starterTemplateName)
	if err := os.Mkdir(path, 0o755); err != nil {
		t.Fatal(err)
	}
	stubResolveRoot(t, root, nil)

	got, created, err := InstallStarter()
	if err != nil {
		t.Fatalf("InstallStarter() error = %v", err)
	}
	if created {
		t.Fatal("InstallStarter() created = true, want false")
	}
	want := Template{Name: starterTemplateName, Path: path}
	if got != want {
		t.Fatalf("InstallStarter() = %#v, want %#v", got, want)
	}
}

func TestInstallStarterRejectsFileAtTemplatePath(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, starterTemplateName)
	if err := os.WriteFile(path, []byte("not a template directory"), 0o644); err != nil {
		t.Fatal(err)
	}
	stubResolveRoot(t, root, nil)

	_, created, err := InstallStarter()
	if created {
		t.Fatal("InstallStarter() created = true, want false")
	}
	if err == nil || err.Error() != "starter template path is not a directory" {
		t.Fatalf("InstallStarter() error = %v", err)
	}
}

func TestUninstallStarter(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, starterTemplateName)
	if err := os.MkdirAll(filepath.Join(path, "src"), 0o755); err != nil {
		t.Fatal(err)
	}
	stubResolveRoot(t, root, nil)

	_, removed, err := UninstallStarter()
	if err != nil {
		t.Fatalf("UninstallStarter() error = %v", err)
	}
	if !removed {
		t.Fatal("UninstallStarter() removed = false, want true")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("starter path still exists, stat error = %v", err)
	}
}

func TestUninstallStarterMissingTemplate(t *testing.T) {
	root := t.TempDir()
	stubResolveRoot(t, root, nil)

	got, removed, err := UninstallStarter()
	if err != nil {
		t.Fatalf("UninstallStarter() error = %v", err)
	}
	if removed {
		t.Fatal("UninstallStarter() removed = true, want false")
	}
	want := Template{Name: starterTemplateName, Path: filepath.Join(root, starterTemplateName)}
	if got != want {
		t.Fatalf("UninstallStarter() = %#v, want %#v", got, want)
	}
}
