package cli

import (
	"fmt"
	"io"

	"github.com/kulikov-andrej/mkproj/internal/core"
)

func Run(
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	if opts.help {
		showHelp(stdout)
		return nil
	}

	if opts.list {
		items, err := core.ListTemplates()
		if err != nil {
			return err
		}
		if len(items) == 0 {
			fmt.Fprintln(stderr, "No templates found.")
			return nil
		}

		for _, template := range items {
			fmt.Fprintln(stdout, template.Name)
		}

		return nil
	}

	if opts.template == "" {
		return fmt.Errorf("template is required")
	}

	if opts.target == "" {
		return fmt.Errorf("target path is required")
	}

	fmt.Fprintf(
		stdout,
		"Creating project %q using template %q...\n",
		opts.target,
		opts.template,
	)

	project, err := core.CreateProject(
		opts.template,
		opts.target,
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
		return core.OpenProject(project)
	}

	return nil
}
