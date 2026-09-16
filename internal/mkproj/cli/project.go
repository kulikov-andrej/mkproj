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

func runWorkflow(
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	if len(args) > 1 {
		return fmt.Errorf("unexpected argument %q", args[1])
	}

	proj, projectErr := getCurrentProject()

	if len(args) == 0 {
		showWorkflowHelp(stdout)

		if projectErr != nil {
			return nil
		}

		commands, err := listWorkflowCommands(
			proj,
			stdin,
			stdout,
			stderr,
		)
		if err != nil {
			return err
		}

		if len(commands) != 0 {
			fmt.Fprintln(stdout)
			showWorkflowCommands(stdout, commands)
		}

		return nil
	}

	if projectErr != nil {
		return projectErr
	}

	return runProjectWorkflow(
		proj,
		args[0],
		stdin,
		stdout,
		stderr,
	)
}

func showWorkflowHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, "Usage:")
	fmt.Fprintln(stdout, "  mkproj run <command>")
}

func showRunHelp(stdout io.Writer) {
	showWorkflowHelp(stdout)
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Run project workflow commands defined in `.mkproj/workflow.star`.")
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, `Without a command, lists available workflow commands.

Examples:
  mkproj run
  mkproj run build
  mkproj run test`)
}

func showWorkflowCommands(
	stdout io.Writer,
	commands []string,
) {
	fmt.Fprintln(stdout, "Commands:")

	for _, command := range commands {
		fmt.Fprintf(stdout, "  %s\n", command)
	}
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
