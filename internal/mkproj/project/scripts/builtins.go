package scripts

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kulikov-andrej/mkproj/internal/libfs"
	"github.com/kulikov-andrej/mkproj/internal/libstarlark"

	"go.starlark.net/starlark"
)

func (ctx *Context) run(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	if len(kwargs) != 0 {
		return nil, fmt.Errorf(
			"%s: keyword arguments are not supported",
			builtin.Name(),
		)
	}

	if len(args) == 0 {
		return nil, fmt.Errorf(
			"%s: expected a program name",
			builtin.Name(),
		)
	}

	command, err := libstarlark.AsString(args[0])
	if err != nil {
		return nil, fmt.Errorf(
			"%s: program: %w",
			builtin.Name(),
			err,
		)
	}

	commandArgs := make([]string, 0, len(args)-1)
	for i, value := range args[1:] {
		arg, err := libstarlark.AsString(value)
		if err != nil {
			return nil, fmt.Errorf(
				"%s: argument %d: %w",
				builtin.Name(),
				i+1,
				err,
			)
		}

		commandArgs = append(commandArgs, arg)
	}

	cmd := exec.Command(command, commandArgs...)
	cmd.Dir = ctx.ProjectPath
	cmd.Stdin = ctx.Streams.In
	cmd.Stdout = ctx.Streams.Out
	cmd.Stderr = ctx.Streams.Err

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			builtin.Name(),
			command,
			err,
		)
	}

	return starlark.None, nil
}

func (ctx *Context) replace(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var path string
	var old string
	var newValue string

	if err := starlark.UnpackArgs(
		builtin.Name(),
		args,
		kwargs,
		"path",
		&path,
		"old",
		&old,
		"new",
		&newValue,
	); err != nil {
		return nil, err
	}

	fullPath, err := libfs.ResolveInside(ctx.ProjectPath, path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf(
			"%s: read %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	text := string(content)
	count := strings.Count(text, old)
	if count == 0 {
		return starlark.MakeInt(0), nil
	}

	text = strings.ReplaceAll(text, old, newValue)
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(fullPath, []byte(text), info.Mode().Perm()); err != nil {
		return nil, fmt.Errorf(
			"%s: write %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.MakeInt(count), nil
}

func (ctx *Context) write(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var path string
	var content string

	if err := starlark.UnpackArgs(
		builtin.Name(),
		args,
		kwargs,
		"path",
		&path,
		"content",
		&content,
	); err != nil {
		return nil, err
	}

	fullPath, err := libfs.ResolveInside(ctx.ProjectPath, path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return nil, fmt.Errorf(
			"%s: create parent directory: %w",
			builtin.Name(),
			err,
		)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		return nil, fmt.Errorf(
			"%s: write %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.None, nil
}

func (ctx *Context) mkdir(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var path string

	if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "path", &path); err != nil {
		return nil, err
	}

	fullPath, err := libfs.ResolveInside(ctx.ProjectPath, path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(fullPath, 0o755); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.None, nil
}

func (ctx *Context) remove(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var path string

	if err := starlark.UnpackArgs(builtin.Name(), args, kwargs, "path", &path); err != nil {
		return nil, err
	}

	fullPath, err := libfs.ResolveInside(ctx.ProjectPath, path)
	if err != nil {
		return nil, err
	}

	if fullPath == ctx.ProjectPath {
		return nil, fmt.Errorf(
			"%s: refusing to remove project root",
			builtin.Name(),
		)
	}

	if err := os.RemoveAll(fullPath); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.None, nil
}
