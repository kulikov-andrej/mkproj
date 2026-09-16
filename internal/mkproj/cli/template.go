package cli

import (
	"fmt"
	"io"
)

func runTemplate(
	args []string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	if len(args) == 0 {
		showTemplateHelp(stdout)
		return nil
	}

	switch args[0] {
	case "list":
		if len(args) > 1 {
			return fmt.Errorf("unexpected argument %q", args[1])
		}

		return runTemplateList(stdout, stderr)

	default:
		return fmt.Errorf("unknown template command %q", args[0])
	}
}

func runTemplateList(
	stdout io.Writer,
	stderr io.Writer,
) error {
	tmpls, err := listTemplates()
	if err != nil {
		return err
	}

	if len(tmpls) == 0 {
		fmt.Fprintln(stderr, "No templates found.")
		return nil
	}

	for _, tmpl := range tmpls {
		fmt.Fprintln(stdout, tmpl.Name)
	}

	return nil
}

func showTemplateHelp(stdout io.Writer) {
	fmt.Fprintln(stdout, `Usage:
  mkproj template <command> [<args>]

Commands:
  list                       List available templates`)
}
