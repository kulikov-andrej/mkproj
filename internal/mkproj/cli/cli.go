package cli

import (
	"github.com/kulikov-andrej/mkproj/internal/libio"
)

func Run(
	args []string,
	streams libio.IO,
) error {
	if len(args) == 0 {
		showHelp(streams)
		return nil
	}

	switch args[0] {
	case "help":
		return runHelp(args[1:], streams)

	case "version":
		return runVersion(streams)

	case "template":
		return runTemplate(args[1:], streams)

	default:
		return runProject(args, streams)
	}
}
