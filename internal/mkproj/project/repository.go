package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kulikov-andrej/mkproj/internal/libfs"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func createProject(
	template templates.Template,
	target string,
) (Project, error) {
	path, err := prepareTarget(target)
	if err != nil {
		return Project{}, err
	}

	if err := libfs.CopyTree(
		template.Path,
		path,
	); err != nil {
		return Project{}, fmt.Errorf(
			"copy template: %w",
			err,
		)
	}

	return Project{
		Name: filepath.Base(path),
		Path: path,
	}, nil
}

func prepareTarget(target string) (string, error) {
	if target == "" {
		return "", fmt.Errorf(
			"target path cannot be empty",
		)
	}

	path, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf(
			"resolve target path: %w",
			err,
		)
	}

	path = filepath.Clean(path)

	info, err := os.Stat(path)

	switch {
	case os.IsNotExist(err):
		return path, nil

	case err != nil:
		return "", fmt.Errorf(
			"access target: %w",
			err,
		)

	case !info.IsDir():
		return "", fmt.Errorf(
			"target is not a directory: %s",
			path,
		)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf(
			"read target directory: %w",
			err,
		)
	}

	if len(entries) != 0 {
		return "", fmt.Errorf(
			"target directory is not empty: %s",
			path,
		)
	}

	return path, nil
}
