package project

import (
	"io"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/model"
)

type Project = model.Project

func Create(
	templateName string,
	target string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) (Project, error) {
	tmpl, err := getTemplate(templateName)
	if err != nil {
		return Project{}, err
	}

	proj, err := createProject(tmpl, target)
	if err != nil {
		return Project{}, err
	}

	setupErr := runSetup(
		proj,
		tmpl,
		stdin,
		stdout,
		stderr,
	)

	cleanupErr := cleanupMetadata(proj)

	if setupErr != nil {
		return Project{}, setupErr
	}

	if cleanupErr != nil {
		return Project{}, cleanupErr
	}

	return proj, nil
}

func Open(proj Project) error {
	return openEditor(proj.Path)
}
