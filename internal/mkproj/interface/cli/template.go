package cli

import (
	"fmt"
	"io"
)

func runTemplateList(
	stdout io.Writer,
	stderr io.Writer,
) error {
	items, err := listTemplates()
	if err != nil {
		return err
	}

	if len(items) == 0 {
		fmt.Fprintln(stderr, "No templates found.")
		return nil
	}

	for _, template := range items {
		fmt.Fprintln(stdout, template.Name)
	}

	return nil
}
