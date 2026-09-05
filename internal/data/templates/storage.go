package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func defaultRoot() (string, error) {
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

func List() ([]Template, error) {
	root, err := resolveRoot()
	if err != nil {
		return []Template{}, err
	}

	entries, err := os.ReadDir(root)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf(
			"read templates directory: %w",
			err,
		)
	}

	result := make(
		[]Template,
		0,
		len(entries),
	)

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		result = append(
			result,
			Template{
				Name: entry.Name(),
				Path: filepath.Join(
					root,
					entry.Name(),
				),
			},
		)
	}

	sort.Slice(
		result,
		func(i, j int) bool {
			return result[i].Name < result[j].Name
		},
	)

	return result, nil
}

func Get(name string) (Template, error) {
	if err := validateName(name); err != nil {
		return Template{}, err
	}

	root, err := resolveRoot()
	if err != nil {
		return Template{}, err
	}

	path := filepath.Join(
		root,
		name,
	)

	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return Template{}, fmt.Errorf(
			"template %q not found",
			name,
		)
	}

	if err != nil {
		return Template{}, fmt.Errorf(
			"access template %q: %w",
			name,
			err,
		)
	}

	if !info.IsDir() {
		return Template{}, fmt.Errorf(
			"template %q is not a directory",
			name,
		)
	}

	return Template{
		Name: name,
		Path: path,
	}, nil
}

func validateName(name string) error {
	if name == "" {
		return fmt.Errorf(
			"template name cannot be empty",
		)
	}

	if name == "." ||
		name == ".." ||
		strings.ContainsAny(name, `/\`) {
		return fmt.Errorf(
			"invalid template name: %q",
			name,
		)
	}

	return nil
}
