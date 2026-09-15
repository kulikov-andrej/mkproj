package cli

import (
	"fmt"
	"io"
)

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
