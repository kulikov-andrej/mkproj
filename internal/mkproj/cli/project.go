package cli

import (
	"fmt"
	"io"
)

func runProject(
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	opts, err := parseProjectArgs(args)
	if err != nil {
		return err
	}

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

	proj, err := createProject(
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
		proj.Name,
		proj.Path,
	)

	if opts.open {
		fmt.Fprintln(stdout, "Opening Code...")

		if err := openProject(proj); err != nil {
			return err
		}
	}

	fmt.Fprintln(stdout, "Have a nice day :)")

	return nil
}

func showProjectHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, `Usage:
  mkproj [<path>] -t <template> [--open]

Arguments:
  path                       Project path (default: current directory)

Options:
  -t, --template <template>  Template name to use
  -o, --open                 Open project in Code`)
}
