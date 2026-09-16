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

func GetCurrentProject() (projectmodel.Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return projectmodel.Project{}, err
	}

	for {
		metadataPath := filepath.Join(dir, ".mkproj")
		info, err := os.Stat(metadataPath)

		if err == nil && info.IsDir() {
			return projectmodel.Project{
				Name: filepath.Base(dir),
				Path: dir,
			}, nil
		}

		if err != nil && !os.IsNotExist(err) {
			return projectmodel.Project{}, fmt.Errorf(
				"access project metadata: %w",
				err,
			)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
	}

	return projectmodel.Project{}, fmt.Errorf("not inside an mkproj project")
}

func CleanupMetadata(proj projectmodel.Project) error {
	metadataPath := filepath.Join(
		proj.Path,
		".mkproj",
	)

	setupPath := filepath.Join(
		metadataPath,
		"setup.star",
	)

	if err := os.RemoveAll(setupPath); err != nil {
		return fmt.Errorf(
			"remove setup metadata: %w",
			err,
		)
	}

	entries, err := os.ReadDir(metadataPath)
	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf(
			"read project metadata: %w",
			err,
		)
	}

	if len(entries) == 0 {
		if err := os.Remove(metadataPath); err != nil {
			return fmt.Errorf(
				"remove empty project metadata: %w",
				err,
			)
		}
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
