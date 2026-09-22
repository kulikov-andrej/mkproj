package libfs

import (
	"path/filepath"
	"testing"
)

func TestResolveInside(t *testing.T) {
	root := t.TempDir()

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name: "file",
			path: "file.txt",
		},
		{
			name: "nested file",
			path: filepath.Join(
				"src",
				"main.cpp",
			),
		},
		{
			name:    "empty path",
			path:    "",
			wantErr: true,
		},
		{
			name:    "parent",
			path:    "..",
			wantErr: true,
		},
		{
			name: "escape root",
			path: filepath.Join(
				"..",
				"outside.txt",
			),
			wantErr: true,
		},
		{
			name:    "absolute path",
			path:    filepath.Join(root, "outside.txt"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResolveInside(
				root,
				tt.path,
			)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}

				return
			}

			if err != nil {
				t.Fatal(err)
			}

			want := filepath.Join(
				root,
				tt.path,
			)

			if got != filepath.Clean(want) {
				t.Fatalf(
					"expected %q, got %q",
					want,
					got,
				)
			}
		})
	}
}
