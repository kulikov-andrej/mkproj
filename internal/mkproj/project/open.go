package project

import (
	"fmt"
	"os/exec"
)

func Open(proj Project) error {
	code, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf(
			"code executable not found: %w",
			err,
		)
	}

	if err := exec.Command(code, proj.Path).Run(); err != nil {
		return fmt.Errorf(
			"open project in Code: %w",
			err,
		)
	}

	return nil
}
