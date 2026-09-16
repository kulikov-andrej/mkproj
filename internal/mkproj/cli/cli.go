package cli

import (
	"io"
)

func Run(
	args []string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	if len(args) == 0 {
		showHelp(stdout)
		return nil
	}

	switch args[0] {
	case "help":
		return runHelp(args[1:], stdout)

	case "version":
		return runVersion(stdout)

	case "template":
		return runTemplate(args[1:], stdout, stderr)

	case "run":
		return runWorkflow(args[1:], stdin, stdout, stderr)

	default:
		return runProject(args, stdin, stdout, stderr)
	}
}
