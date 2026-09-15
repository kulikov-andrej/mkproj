package setup

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	projectmodel "github.com/kulikov-andrej/mkproj/internal/mkproj/project/model"
	templatemodel "github.com/kulikov-andrej/mkproj/internal/mkproj/templates/model"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

type setupContext struct {
	proj   projectmodel.Project
	tmpl   templatemodel.Template
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func Run(
	proj projectmodel.Project,
	tmpl templatemodel.Template,
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
		proj:   proj,
		tmpl:   tmpl,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
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

func (ctx *setupContext) globals() starlark.StringDict {
	return starlark.StringDict{
		"project": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(ctx.proj.Name),
				"path": starlark.String(ctx.proj.Path),
			},
		),

		"template": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(ctx.tmpl.Name),
			},
		),

		"run": starlark.NewBuiltin(
			"run",
			ctx.run,
		),

		"replace": starlark.NewBuiltin(
			"replace",
			ctx.replace,
		),

		"write": starlark.NewBuiltin(
			"write",
			ctx.write,
		),

		"mkdir": starlark.NewBuiltin(
			"mkdir",
			ctx.mkdir,
		),

		"remove": starlark.NewBuiltin(
			"remove",
			ctx.remove,
		),
	}
}
