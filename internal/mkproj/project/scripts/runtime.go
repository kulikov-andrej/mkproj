package scripts

import (
	"fmt"
	"os"
	"sort"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/libstarlark"

	"go.starlark.net/starlark"
)

func ExecFile(
	path string,
	ctx Context,
) error {
	_, err := LoadFile(path, ctx)
	return err
}

func LoadFile(
	path string,
	ctx Context,
) (starlark.StringDict, error) {
	return libstarlark.ExecFile(path, buildAPI(ctx), ctx.Streams)
}

func Call(
	fn starlark.Callable,
	ctx Context,
) error {
	_, err := starlark.Call(
		newThread(ctx.Streams),
		fn,
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("script failed: %w", err)
	}

	return nil
}

func Commands(
	path string,
	ctx Context,
) ([]string, error) {
	globals, err := LoadFile(path, ctx)
	if os.IsNotExist(err) {
		return []string{}, nil
	}
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

func Run(
	path string,
	command string,
	ctx Context,
) error {
	globals, err := LoadFile(path, ctx)
	if os.IsNotExist(err) {
		globals = starlark.StringDict{}
	}
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	value, ok := globals[command]
	if !ok {
		return fmt.Errorf("workflow command %q not found", command)
	}

	fn, ok := value.(starlark.Callable)
	if !ok {
		return fmt.Errorf("workflow command %q is not callable", command)
	}

	return Call(fn, ctx)
}

func buildAPI(ctx Context) starlark.StringDict {
	return libstarlark.NewAPI().
		Struct("project", map[string]string{
			"name": ctx.ProjectName,
			"path": ctx.ProjectPath,
		}).
		Struct("template", map[string]string{
			"name": ctx.TemplateName,
		}).
		Function("run", ctx.run).
		Function("replace", ctx.replace).
		Function("write", ctx.write).
		Function("mkdir", ctx.mkdir).
		Function("remove", ctx.remove).
		Build()
}

func newThread(streams libio.IO) *starlark.Thread {
	return &starlark.Thread{
		Name: "mkproj",
		Print: func(
			_ *starlark.Thread,
			message string,
		) {
			fmt.Fprintln(streams.Out, message)
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
