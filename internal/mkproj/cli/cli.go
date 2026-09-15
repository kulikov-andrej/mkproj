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
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	switch {
	case len(args) == 0 || opts.help:
		showHelp(stdout)
		return nil

	case opts.version:
		return runVersion(stdout)

	case opts.list:
		return runTemplateList(stdout, stderr)

	default:
		return runProject(
			opts,
			stdin,
			stdout,
			stderr,
		)
	}
}
