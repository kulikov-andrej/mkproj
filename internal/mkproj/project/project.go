package project

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/model"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/scripts"
	"go.starlark.net/starlark"
)

type Project = model.Project

func GetCurrent() (Project, error) {
	return getCurrentProject()
}

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

	setupPath := filepath.Join(
		proj.Path,
		".mkproj",
		"setup.star",
	)

	_, setupErr := loadScript(
		setupPath,
		scriptContext(
			proj,
			tmpl.Name,
			stdin,
			stdout,
			stderr,
		),
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

func WorkflowCommands(
	proj Project,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) ([]string, error) {
	globals, err := loadWorkflow(
		proj,
		scriptContext(
			proj,
			"",
			stdin,
			stdout,
			stderr,
		),
	)
	if err != nil {
		return nil, err
	}

	commands := make([]string, 0, len(globals))

	for name, value := range globals {
		if _, ok := value.(starlark.Callable); ok {
			commands = append(commands, name)
		}
	}

	sort.Strings(commands)

	return commands, nil
}

func RunWorkflow(
	proj Project,
	command string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	ctx := scriptContext(
		proj,
		"",
		stdin,
		stdout,
		stderr,
	)

	globals, err := loadWorkflow(proj, ctx)
	if err != nil {
		return err
	}

	value, ok := globals[command]
	if !ok {
		return fmt.Errorf("workflow command %q not found", command)
	}

	fn, ok := value.(starlark.Callable)
	if !ok {
		return fmt.Errorf(
			"workflow command %q is not callable",
			command,
		)
	}

	return callScript(fn, ctx)
}

func loadWorkflow(
	proj Project,
	ctx scripts.Context,
) (starlark.StringDict, error) {
	return loadScript(
		filepath.Join(
			proj.Path,
			".mkproj",
			"workflow.star",
		),
		ctx,
	)
}

func scriptContext(
	proj Project,
	templateName string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) scripts.Context {
	return scripts.Context{
		ProjectName:  proj.Name,
		ProjectPath:  proj.Path,
		TemplateName: templateName,
		Stdin:        stdin,
		Stdout:       stdout,
		Stderr:       stderr,
	}
}

func Open(proj Project) error {
	return openEditor(proj.Path)
}
