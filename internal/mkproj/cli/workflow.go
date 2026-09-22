package cli

import (
	"fmt"

	"github.com/kulikov-andrej/mkproj/internal/libio"
)

func runWorkflow(
	args []string,
	streams libio.IO,
) error {
	if len(args) > 1 {
		return fmt.Errorf("unexpected argument %q", args[1])
	}

	proj, projectErr := getProject()
	if len(args) == 0 {
		showWorkflowHelp(streams)

		if projectErr != nil {
			return nil
		}

		commands, err := listWorkflow(proj, streams)
		if err != nil {
			return err
		}

		if len(commands) != 0 {
			fmt.Fprintln(streams.Out)
			showWorkflowCommands(streams, commands)
		}

		return nil
	}

	if projectErr != nil {
		return projectErr
	}

	return runWorkflowFn(proj, args[0], streams)
}

func showWorkflowHelp(streams libio.IO) {
	fmt.Fprintln(streams.Out, "Usage:")
	fmt.Fprintln(streams.Out, "  mkproj run <command>")
}

func showWorkflowCommands(
	streams libio.IO,
	commands []string,
) {
	fmt.Fprintln(streams.Out, "Commands:")
	for _, command := range commands {
		fmt.Fprintf(streams.Out, "  %s\n", command)
	}
}
