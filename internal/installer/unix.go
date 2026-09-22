//go:build linux || darwin

package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func install(source string) (installPaths, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return installPaths{}, fmt.Errorf("get home directory: %w", err)
	}

	installDir := filepath.Join(home, ".local", "bin")
	installPath := filepath.Join(installDir, "mkproj")
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return installPaths{}, fmt.Errorf("create install directory: %w", err)
	}
	if err := copyFile(source, installPath); err != nil {
		return installPaths{}, fmt.Errorf("install binary: %w", err)
	}
	if err := os.Chmod(installPath, 0o755); err != nil {
		return installPaths{}, fmt.Errorf("make binary executable: %w", err)
	}

	configHome, err := os.UserConfigDir()
	if err != nil {
		return installPaths{}, fmt.Errorf(
			"get user config directory: %w",
			err,
		)
	}
	templates := templateDir(configHome)
	created, err := ensureDirectory(templates)
	if err != nil {
		return installPaths{}, err
	}

	return installPaths{Binary: installPath, Templates: templates, CreatedTemplates: created}, nil
}

func ensureDirectory(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return false, fmt.Errorf("template path is not a directory: %s", path)
		}
		return false, nil
	}
	if !os.IsNotExist(err) {
		return false, fmt.Errorf("access template directory: %w", err)
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return false, fmt.Errorf("create template directory: %w", err)
	}
	return true, nil
}

func reportPath(paths installPaths, output io.Writer) error {
	if os.Getenv("PATH") == "" {
		return nil
	}
	for _, entry := range filepath.SplitList(os.Getenv("PATH")) {
		if filepath.Clean(entry) == filepath.Clean(filepath.Dir(paths.Binary)) {
			return nil
		}
	}
	fmt.Fprintln(output)
	fmt.Fprintf(output, "%q is not in PATH. Add it to your shell configuration.\n", filepath.Dir(paths.Binary))
	return nil
}
