package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func resolveTargetPath(target string) (string, error) {
	path, err := filepath.Abs(target)
	if err != nil {
		return "", fmt.Errorf(
			"resolve target path: %w",
			err,
		)
	}

	return filepath.Clean(path), nil
}

func prepareTargetDirectory(path string) error {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.MkdirAll(path, 0o755)
	}

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf(
			"target exists and is not a directory: %s",
			path,
		)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if len(entries) != 0 {
		return fmt.Errorf(
			"target directory is not empty: %s",
			path,
		)
	}

	return nil
}

func copyTemplate(source, target string) error {
	return filepath.WalkDir(
		source,
		func(
			path string,
			entry fs.DirEntry,
			walkErr error,
		) error {
			if walkErr != nil {
				return walkErr
			}

			relative, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}

			if relative == "." {
				return nil
			}

			destination := filepath.Join(
				target,
				relative,
			)

			info, err := entry.Info()
			if err != nil {
				return err
			}

			switch {
			case entry.IsDir():
				return os.MkdirAll(
					destination,
					info.Mode().Perm(),
				)

			case entry.Type()&os.ModeSymlink != 0:
				link, err := os.Readlink(path)
				if err != nil {
					return err
				}

				return os.Symlink(
					link,
					destination,
				)

			case info.Mode().IsRegular():
				return copyFile(
					path,
					destination,
					info.Mode().Perm(),
				)

			default:
				return fmt.Errorf(
					"unsupported file type: %s",
					path,
				)
			}
		},
	)
}

func copyFile(
	source string,
	target string,
	mode fs.FileMode,
) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.OpenFile(
		target,
		os.O_CREATE|
			os.O_TRUNC|
			os.O_WRONLY,
		mode,
	)

	if err != nil {
		return err
	}

	_, copyErr := io.Copy(output, input)
	closeErr := output.Close()

	if copyErr != nil {
		return copyErr
	}

	return closeErr
}
