package libfs

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyTree(source, target string) error {
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

			info, err := entry.Info()
			if err != nil {
				return err
			}

			relative, err := filepath.Rel(source, path)
			if err != nil {
				return err
			}

			if relative == "." {
				if !entry.IsDir() {
					return fmt.Errorf(
						"source is not a directory: %s",
						source,
					)
				}

				return os.MkdirAll(
					target,
					info.Mode().Perm(),
				)
			}

			destination := filepath.Join(
				target,
				relative,
			)

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
