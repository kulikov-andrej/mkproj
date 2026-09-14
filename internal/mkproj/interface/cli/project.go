package cli

import (
	"fmt"
	"io"
)

func runProject(
	opts options,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	if opts.template == "" {
		return fmt.Errorf("template is required")
	}

	target := opts.target
	if target == "" {
		target = "."
	}

	fmt.Fprintf(
		stdout,
		"Creating project %q using template %q...\n",
		target,
		opts.template,
	)

	project, err := createProject(
		opts.template,
		target,
		stdin,
		stdout,
		stderr,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		stdout,
		"Created %q.\n  %s\n",
		project.Name,
		project.Path,
	)

	if opts.open {
		fmt.Fprintln(stdout, "Opening Code...")

		if err := openProject(project); err != nil {
			return err
		}
	}

	fmt.Fprintln(stdout, "Have a nice day :)")

	return nil
}
