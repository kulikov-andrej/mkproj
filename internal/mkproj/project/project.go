package project

import (
	"github.com/kulikov-andrej/mkproj/internal/libio"
)

func Create(
	templateName string,
	projectPath string,
	streams libio.IO,
) (Project, error) {
	template, err := findTemplate(templateName)
	if err != nil {
		return Project{}, err
	}

	proj, err := createProject(template, projectPath)
	if err != nil {
		return Project{}, err
	}
	return proj, nil
}
