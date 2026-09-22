package cli

import (
	"fmt"

	"github.com/kulikov-andrej/mkproj/internal/libio"
)

func runHelp(
	args []string,
	streams libio.IO,
) error {
	if len(args) == 0 {
		showHelp(streams)
		return nil
	}

	topic := args[0]
	switch topic {
	case "project":
		showProjectHelp(streams)
	case "template":
		showTemplateHelp(streams)
	case "run":
		showWorkflowHelp(streams)
	default:
		return fmt.Errorf("unknown topic %q", topic)
	}

	return nil
}

func showHelp(
	streams libio.IO,
) {
	fmt.Fprintln(streams.Out, `Usage:
  mkproj [<path>] -t <template> [--open]
  mkproj <command> [<args>]

Commands:
  help       Show help
  version    Show version
  template   Manage templates
  run        Run project workflow

Help topics:
  project    Project creation
  template   Template commands
  run        Project workflow`)

	fmt.Fprintln(streams.Out)
	fmt.Fprintln(streams.Out, "Run `mkproj help <topic>` for more information.")
}
