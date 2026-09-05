package core

import (
	"io"

	"github.com/kulikov-andrej/mkproj/internal/data/project"
)

func CreateProject(
	templateName string,
	target string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (project.Project, error) {
	template, err := getTemplate(templateName)
	if err != nil {
		return project.Project{}, err
	}

	proj, err := createProject(template, target)
	if err != nil {
		return project.Project{}, err
	}

	hookErr := runHook(
		proj,
		template,
		stdin,
		stdout,
		stderr,
	)

	cleanupErr := cleanupMetadata(proj)

	if hookErr != nil {
		return project.Project{}, hookErr
	}

	if cleanupErr != nil {
		return project.Project{}, cleanupErr
	}

	return proj, nil
}

func OpenProject(proj project.Project) error {
	return openProject(proj.Path)
}
