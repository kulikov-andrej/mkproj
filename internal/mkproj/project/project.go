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

	setupErr := runSetup(proj, template, streams)
	cleanupErr := cleanupSetup(proj)
	if setupErr != nil {
		return Project{}, setupErr
	}
	if cleanupErr != nil {
		return Project{}, cleanupErr
	}

	return proj, nil
}
