package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func ListTemplates() ([]Template, error) {
	root, err := resolveRoot()
	if err != nil {
		return []Template{}, err
	}

	entries, err := readDir(root)

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

const starterTemplateName = "starter"

func InstallStarter() (Template, bool, error) {
	root, err := resolveRoot()
	if err != nil {
		return Template{}, false, err
	}

	template := Template{
		Name: starterTemplateName,
		Path: filepath.Join(root, starterTemplateName),
	}

	info, err := os.Stat(template.Path)
	if err == nil {
		if !info.IsDir() {
			return Template{}, false, fmt.Errorf("starter template path is not a directory")
		}
		return template, false, nil
	}
	if !os.IsNotExist(err) {
		return Template{}, false, fmt.Errorf("access starter template: %w", err)
	}

	if err := os.MkdirAll(filepath.Join(template.Path, ".mkproj"), 0o755); err != nil {
		return Template{}, false, fmt.Errorf("create starter template: %w", err)
	}
	if err := os.MkdirAll(filepath.Join(template.Path, "src"), 0o755); err != nil {
		return Template{}, false, fmt.Errorf("create starter template: %w", err)
	}

	files := map[string]string{
		"README.md":                      "# {{PROJECT_NAME}}\n\nGenerated with the mkproj starter template.\n",
		filepath.Join("src", "main.txt"): "Hello from {{PROJECT_NAME}}!\n",
		filepath.Join(".mkproj", "setup.star"): `replace("README.md", "{{PROJECT_NAME}}", project.name)
replace("src/main.txt", "{{PROJECT_NAME}}", project.name)
`,
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(template.Path, name), []byte(content), 0o644); err != nil {
			_ = os.RemoveAll(template.Path)
			return Template{}, false, fmt.Errorf("write starter template: %w", err)
		}
	}

	return template, true, nil
}

func UninstallStarter() (Template, bool, error) {
	root, err := resolveRoot()
	if err != nil {
		return Template{}, false, err
	}

	template := Template{
		Name: starterTemplateName,
		Path: filepath.Join(root, starterTemplateName),
	}

	if _, err := os.Stat(template.Path); os.IsNotExist(err) {
		return template, false, nil
	} else if err != nil {
		return Template{}, false, fmt.Errorf("access starter template: %w", err)
	}

	if err := os.RemoveAll(template.Path); err != nil {
		return Template{}, false, fmt.Errorf("remove starter template: %w", err)
	}

	return template, true, nil
}

func FindTemplate(
	name string,
) (Template, error) {
	if err := validateTemplateName(name); err != nil {
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

func getDefaultRoot() (string, error) {
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

func validateTemplateName(name string) error {
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
