package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kulikov-andrej/mkproj/internal/libfs"
	projectmodel "github.com/kulikov-andrej/mkproj/internal/mkproj/project/model"
	templatemodel "github.com/kulikov-andrej/mkproj/internal/mkproj/templates/model"
)

func Create(
	tmpl templatemodel.Template,
	target string,
) (projectmodel.Project, error) {
	path, err := prepareTarget(target)
	if err != nil {
		return projectmodel.Project{}, err
	}

	if err := libfs.CopyTree(
		tmpl.Path,
		path,
	); err != nil {
		return projectmodel.Project{}, fmt.Errorf(
			"copy template: %w",
			err,
		)
	}

	return projectmodel.Project{
		Name: filepath.Base(path),
		Path: path,
	}, nil
}

func CleanupMetadata(proj projectmodel.Project) error {
	path := filepath.Join(
		proj.Path,
		".mkproj",
	)

	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf(
			"remove project metadata: %w",
			err,
		)
	}

	return nil
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
