package libfs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyTree(t *testing.T) {
	root := t.TempDir()

	source := filepath.Join(root, "source")
	target := filepath.Join(root, "target")

	if err := os.MkdirAll(
		filepath.Join(source, "src"),
		0o755,
	); err != nil {
		t.Fatal(err)
	}

	files := map[string]string{
		"README.md":    "# example\n",
		".gitignore":   "/build/\n",
		"src/main.cpp": "int main() {}\n",
	}

	for path, content := range files {
		fullPath := filepath.Join(
			source,
			filepath.FromSlash(path),
		)

		if err := os.WriteFile(
			fullPath,
			[]byte(content),
			0o644,
		); err != nil {
			t.Fatal(err)
		}
	}

	if err := CopyTree(source, target); err != nil {
		t.Fatal(err)
	}

	for path, expected := range files {
		fullPath := filepath.Join(
			target,
			filepath.FromSlash(path),
		)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf(
				"read %q: %v",
				fullPath,
				err,
			)
		}

		if string(content) != expected {
			t.Fatalf(
				"%q: expected %q, got %q",
				fullPath,
				expected,
				string(content),
			)
		}
	}
}

func TestCopyTreeRejectsFileSource(t *testing.T) {
	root := t.TempDir()

	source := filepath.Join(root, "source.txt")
	target := filepath.Join(root, "target")

	if err := os.WriteFile(
		source,
		[]byte("hello"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	if err := CopyTree(source, target); err == nil {
		t.Fatal("expected an error")
	}
}
