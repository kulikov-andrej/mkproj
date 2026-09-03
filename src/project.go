package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func (a *app) newProject(
	template string,
	target string,
	open bool,
) error {
	templatePath, err := a.templatePath(template)
	if err != nil {
		return err
	}

	targetPath, err := resolveTargetPath(target)
	if err != nil {
		return err
	}

	if err := prepareTargetDirectory(targetPath); err != nil {
		return err
	}

	projectName := filepath.Base(targetPath)

	fmt.Fprintf(
		a.stdout,
		"Creating %q from template %q...\n",
		projectName,
		template,
	)

	fmt.Fprintf(a.stdout, "  %s\n", targetPath)

	if err := copyTemplate(
		templatePath,
		targetPath,
	); err != nil {
		return fmt.Errorf(
			"copy template: %w",
			err,
		)
	}

	hookErr := a.invokeTemplateHook(
		targetPath,
		template,
		projectName,
	)

	cleanupErr := os.RemoveAll(
		filepath.Join(
			targetPath,
			".mkproj",
		),
	)

	if hookErr != nil {
		if cleanupErr != nil {
			return fmt.Errorf(
				"%w; cleanup template metadata: %v",
				hookErr,
				cleanupErr,
			)
		}

		return hookErr
	}

	if cleanupErr != nil {
		return fmt.Errorf(
			"cleanup template metadata: %w",
			cleanupErr,
		)
	}

	fmt.Fprintf(
		a.stdout,
		"Created %q.\n",
		projectName,
	)

	if open {
		fmt.Fprintln(a.stdout, "Opening Code...")

		if err := a.openProject(targetPath); err != nil {
			return err
		}
	}

	return nil
}

func (a *app) openInCode(path string) error {
	code, err := exec.LookPath("code")
	if err != nil {
		return fmt.Errorf(
			"'code' was not found in PATH",
		)
	}

	cmd := exec.Command(code, path)

	cmd.Stdin = a.stdin
	cmd.Stdout = a.stdout
	cmd.Stderr = a.stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"failed to open Code: %w",
			err,
		)
	}

	return nil
}
