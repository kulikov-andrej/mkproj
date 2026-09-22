package libfs

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCopyTreeCopiesTree(t *testing.T) {
	source := t.TempDir()
	target := filepath.Join(t.TempDir(), "target")

	if err := os.MkdirAll(filepath.Join(source, "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "root.txt"), []byte("root"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(source, "nested", "child.txt"), []byte("child"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := CopyTree(source, target); err != nil {
		t.Fatalf("CopyTree() error = %v", err)
	}

	assertFileContent(t, filepath.Join(target, "root.txt"), "root")
	assertFileContent(t, filepath.Join(target, "nested", "child.txt"), "child")

	info, err := os.Stat(filepath.Join(target, "nested"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatalf("nested path is not a directory")
	}
}

func TestCopyTreeCopiesEmptyDirectory(t *testing.T) {
	source := t.TempDir()
	if err := os.Mkdir(filepath.Join(source, "empty"), 0o755); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(t.TempDir(), "target")
	if err := CopyTree(source, target); err != nil {
		t.Fatalf("CopyTree() error = %v", err)
	}

	info, err := os.Stat(filepath.Join(target, "empty"))
	if err != nil {
		t.Fatal(err)
	}
	if !info.IsDir() {
		t.Fatal("empty directory was not copied as a directory")
	}
}

func TestCopyTreeOverwritesExistingFile(t *testing.T) {
	source := t.TempDir()
	target := t.TempDir()

	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("new"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(target, "file.txt"), []byte("old content that is longer"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := CopyTree(source, target); err != nil {
		t.Fatalf("CopyTree() error = %v", err)
	}

	assertFileContent(t, filepath.Join(target, "file.txt"), "new")
}

func TestCopyTreeRejectsFileSource(t *testing.T) {
	source := filepath.Join(t.TempDir(), "source.txt")
	if err := os.WriteFile(source, []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}

	err := CopyTree(source, filepath.Join(t.TempDir(), "target"))
	if err == nil {
		t.Fatal("CopyTree() error = nil, want error")
	}
}

func TestCopyTreeCopiesSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink creation may require Windows Developer Mode or elevated privileges")
	}

	source := t.TempDir()
	target := filepath.Join(t.TempDir(), "target")

	if err := os.WriteFile(filepath.Join(source, "file.txt"), []byte("data"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("file.txt", filepath.Join(source, "link.txt")); err != nil {
		t.Fatal(err)
	}

	if err := CopyTree(source, target); err != nil {
		t.Fatalf("CopyTree() error = %v", err)
	}

	link, err := os.Readlink(filepath.Join(target, "link.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if link != "file.txt" {
		t.Fatalf("symlink target = %q, want %q", link, "file.txt")
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := string(data); got != want {
		t.Fatalf("file content = %q, want %q", got, want)
	}
}
