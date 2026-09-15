package cli

import (
	"fmt"
	"io"
)

func runHelp(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		showHelp(stdout)
		return nil
	}

	if len(args) > 1 {
		return fmt.Errorf("unexpected argument %q", args[1])
	}

	switch args[0] {
	case "project":
		showProjectHelp(stdout)
		return nil

	case "template":
		showTemplateHelp(stdout)
		return nil

	default:
		return fmt.Errorf("unknown help topic %q", args[0])
	}
}

func showHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, `Usage:
  mkproj [<path>] -t <template> [--open]
  mkproj <command> [<args>]

Commands:
  help       Show help
  version    Show version
  template   Manage templates

Help topics:
  project    Project creation
  template   Template commands`)

	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Run `mkproj help <topic>` for more information.")
}
