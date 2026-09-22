package project

import (
	"fmt"
	"os"
	"path/filepath"
)

func cleanupSetup(proj Project) error {
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
