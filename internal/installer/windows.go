//go:build windows

package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func install(source string) (installPaths, error) {
	localAppData := os.Getenv("LOCALAPPDATA")
	appData := os.Getenv("APPDATA")
	if localAppData == "" || appData == "" {
		return installPaths{}, fmt.Errorf("required Windows environment variables are missing")
	}

	installDir := filepath.Join(localAppData, "Programs", "mkproj")
	installPath := filepath.Join(installDir, "mkproj.exe")
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return installPaths{}, fmt.Errorf("create install directory: %w", err)
	}
	if err := copyFile(source, installPath); err != nil {
		return installPaths{}, fmt.Errorf("install binary: %w", err)
	}

	templates := templateDir(appData)
	created, err := ensureWindowsDirectory(templates)
	if err != nil {
		return installPaths{}, err
	}
	if err := addUserPath(installDir); err != nil {
		return installPaths{}, err
	}

	return installPaths{Binary: installPath, Templates: templates, CreatedTemplates: created}, nil
}

func ensureWindowsDirectory(path string) (bool, error) {
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

func addUserPath(directory string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("open user environment: %w", err)
	}
	defer key.Close()

	value, _, err := key.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("read user PATH: %w", err)
	}
	for _, entry := range filepath.SplitList(value) {
		if filepath.Clean(entry) == filepath.Clean(directory) {
			return nil
		}
	}
	if value != "" {
		value += ";"
	}
	return key.SetStringValue("Path", value+directory)
}

func reportPath(paths installPaths, output io.Writer) error {
	value := os.Getenv("PATH")
	for _, entry := range filepath.SplitList(value) {
		if filepath.Clean(entry) == filepath.Clean(filepath.Dir(paths.Binary)) {
			return nil
		}
	}
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Open a new terminal to use 'mkproj' from PATH.")
	return nil
}
