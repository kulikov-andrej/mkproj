package editor

import (
	"fmt"
	"os/exec"
)

func Open(path string) error {
	code, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf(
			"code executable not found: %w",
			err,
		)
	}

	if err := exec.Command(code, path).Run(); err != nil {
		return fmt.Errorf(
			"open project in Code: %w",
			err,
		)
	}

	return nil
}
