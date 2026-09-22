package libfs

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ResolveInside(root, path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if filepath.IsAbs(path) {
		return "", fmt.Errorf("path must be relative: %s", path)
	}

	fullPath := filepath.Clean(
		filepath.Join(root, path),
	)

	relative, err := filepath.Rel(root, fullPath)
	if err != nil {
		return "", err
	}

	if relative == ".." ||
		strings.HasPrefix(
			relative,
			".."+string(filepath.Separator),
		) {
		return "", fmt.Errorf(
			"path escapes root directory: %s",
			path,
		)
	}

	return fullPath, nil
}
