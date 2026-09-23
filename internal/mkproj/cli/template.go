package cli

import (
	"fmt"

	"github.com/kulikov-andrej/mkproj/internal/libio"
)

func runTemplate(
	args []string,
	streams libio.IO,
) error {
	if len(args) == 0 {
		showTemplateHelp(streams)
		return nil
	}

	switch args[0] {
	case "list":
		if len(args) > 1 {
			return fmt.Errorf("unexpected argument %q", args[1])
		}

		return runTemplateList(streams)

	case "init":
		if len(args) > 1 {
			return fmt.Errorf("unexpected argument %q", args[1])
		}
		return runTemplateInit(streams)

	case "deinit":
		if len(args) > 1 {
			return fmt.Errorf("unexpected argument %q", args[1])
		}
		return runTemplateDeinit(streams)

	default:
		return fmt.Errorf("unknown template command %q", args[0])
	}
}

func runTemplateList(
	streams libio.IO,
) error {
	template, err := listTemplates()
	if err != nil {
		return err
	}

	if len(template) == 0 {
		fmt.Fprintln(streams.Err, "No templates found.")
		return nil
	}

	for _, tmpl := range template {
		fmt.Fprintln(streams.Out, tmpl.Name)
	}

	return nil
}

func runTemplateInit(
	streams libio.IO,
) error {
	starterTemplate, created, err := installStarterTemplate()
	if err != nil {
		return err
	}
	if !created {
		fmt.Fprintf(streams.Out, "Starter template is already installed: \n  %s\n", starterTemplate.Path)
		return nil
	}

	fmt.Fprintf(streams.Out, "Created starter template: \n  %s\n", starterTemplate.Path)
	return nil
}

func runTemplateDeinit(
	streams libio.IO,
) error {
	starterTemplate, removed, err := uninstallStarterTemplate()
	if err != nil {
		return err
	}
	if !removed {
		fmt.Fprintln(streams.Out, "Starter template is not installed.")
		return nil
	}

	fmt.Fprintf(streams.Out, "Removed starter template: \n  %s\n", starterTemplate.Path)
	return nil
}

func showTemplateHelp(streams libio.IO) {
	fmt.Fprintln(streams.Out, `Usage:
  mkproj template <command> [<args>]

Commands:
  list                       List available templates
  init                       Install starter template for trying out mkproj
  deinit                     Uninstall starter template`)
}
