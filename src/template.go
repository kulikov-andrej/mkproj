package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func defaultTemplateRoot() (string, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf(
			"get user config directory: %w",
			err,
		)
	}

	return filepath.Join(
		configRoot,
		"mkproj",
		"templates",
	), nil
}

func (a *app) showTemplates() error {
	entries, err := os.ReadDir(a.templateRoot)

	if os.IsNotExist(err) {
		fmt.Fprintln(a.stdout, "No templates found.")
		fmt.Fprintf(
			a.stdout,
			"Template directory: %s\n",
			a.templateRoot,
		)
		return nil
	}

	if err != nil {
		return err
	}

	var names []string

	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}

	if len(names) == 0 {
		fmt.Fprintln(a.stdout, "No templates found.")
		return nil
	}

	sort.Strings(names)

	for _, name := range names {
		fmt.Fprintln(a.stdout, name)
	}

	return nil
}

func (a *app) templatePath(name string) (string, error) {
	if err := validateTemplateName(name); err != nil {
		return "", err
	}

	path := filepath.Join(a.templateRoot, name)

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return "", fmt.Errorf(
			"template %q not found",
			name,
		)
	}

	if err != nil {
		return "", err
	}

	if !info.IsDir() {
		return "", fmt.Errorf(
			"template %q is not a directory",
			name,
		)
	}

	return path, nil
}

func validateTemplateName(name string) error {
	if name == "" ||
		name == "." ||
		name == ".." ||
		strings.ContainsAny(name, `/\`) {

		return fmt.Errorf(
			"invalid template name: %q",
			name,
		)
	}

	return nil
}
