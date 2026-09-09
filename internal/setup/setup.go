package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kulikov-andrej/mkproj/internal/data/project"
	"github.com/kulikov-andrej/mkproj/internal/data/templates"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

type setupContext struct {
	projectName  string
	projectPath  string
	templateName string
	stdin        io.Reader
	stdout       io.Writer
	stderr       io.Writer
}

func Run(
	proj project.Project,
	template templates.Template,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) error {
	setupPath := filepath.Join(
		proj.Path,
		".mkproj",
		"setup.star",
	)

	info, err := os.Stat(setupPath)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("stat template setup: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf(
			"template setup is a directory: %s",
			setupPath,
		)
	}

	ctx := &setupContext{
		projectName:  proj.Name,
		projectPath:  proj.Path,
		templateName: template.Name,
		stdin:        stdin,
		stdout:       stdout,
		stderr:       stderr,
	}

	thread := &starlark.Thread{
		Name: "mkproj setup",

		Print: func(
			_ *starlark.Thread,
			message string,
		) {
			fmt.Fprintln(stdout, message)
		},

		Load: func(
			_ *starlark.Thread,
			module string,
		) (starlark.StringDict, error) {
			return nil, fmt.Errorf(
				"load is not supported: %s",
				module,
			)
		},
	}

	options := &syntax.FileOptions{
		TopLevelControl: true,
	}

	_, err = starlark.ExecFileOptions(
		options,
		thread,
		setupPath,
		nil,
		ctx.globals(),
	)

	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf(
				"template setup failed:\n%s",
				evalErr.Backtrace(),
			)
		}

		return fmt.Errorf(
			"template setup failed: %w",
			err,
		)
	}

	return nil
}

func (h *setupContext) globals() starlark.StringDict {
	return starlark.StringDict{
		"project": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(h.projectName),
				"path": starlark.String(h.projectPath),
			},
		),

		"template": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(h.templateName),
			},
		),

		"run": starlark.NewBuiltin(
			"run",
			h.run,
		),

		"replace": starlark.NewBuiltin(
			"replace",
			h.replace,
		),

		"write": starlark.NewBuiltin(
			"write",
			h.write,
		),

		"mkdir": starlark.NewBuiltin(
			"mkdir",
			h.mkdir,
		),

		"remove": starlark.NewBuiltin(
			"remove",
			h.remove,
		),
	}
}
