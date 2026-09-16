package scripts

import (
	"fmt"
	"os"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

func Load(
	path string,
	ctx Context,
) (starlark.StringDict, error) {
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return starlark.StringDict{}, nil
	}

	if err != nil {
		return starlark.StringDict{}, fmt.Errorf("stat script: %w", err)
	}

	if info.IsDir() {
		return starlark.StringDict{}, fmt.Errorf(
			"script is a directory: %s",
			path,
		)
	}

	globals, err := starlark.ExecFileOptions(
		&syntax.FileOptions{
			TopLevelControl: true,
		},
		newThread(ctx),
		path,
		nil,
		ctx.predeclared(),
	)
	if err != nil {
		return starlark.StringDict{}, scriptError(err)
	}

	return globals, nil
}

func Call(
	fn starlark.Callable,
	ctx Context,
) error {
	_, err := starlark.Call(
		newThread(ctx),
		fn,
		nil,
		nil,
	)
	if err != nil {
		return scriptError(err)
	}

	return nil
}

func newThread(ctx Context) *starlark.Thread {
	return &starlark.Thread{
		Name: "mkproj",
		Print: func(
			_ *starlark.Thread,
			msg string,
		) {
			fmt.Fprintln(ctx.Stdout, msg)
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
}

func scriptError(err error) error {
	if evalErr, ok := err.(*starlark.EvalError); ok {
		return fmt.Errorf(
			"script failed:\n%s",
			evalErr.Backtrace(),
		)
	}

	return fmt.Errorf("script failed: %w", err)
}

func (ctx Context) predeclared() starlark.StringDict {
	values := starlark.StringDict{
		"project": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(ctx.ProjectName),
				"path": starlark.String(ctx.ProjectPath),
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

	if ctx.TemplateName != "" {
		values["template"] = starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(ctx.TemplateName),
			},
		)
	}

	return values
}
