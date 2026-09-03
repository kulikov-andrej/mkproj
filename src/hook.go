package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/syntax"
)

type hookContext struct {
	app         *app
	projectName string
	projectPath string
	template    string
}

func (a *app) invokeTemplateHook(
	projectPath string,
	template string,
	projectName string,
) error {
	hookPath := filepath.Join(
		projectPath,
		".mkproj",
		"hook.star",
	)

	info, err := os.Stat(hookPath)

	if os.IsNotExist(err) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("stat template hook: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf(
			"template hook is a directory: %s",
			hookPath,
		)
	}

	ctx := &hookContext{
		app:         a,
		projectName: projectName,
		projectPath: projectPath,
		template:    template,
	}

	thread := &starlark.Thread{
		Name: "mkproj hook",

		Print: func(
			_ *starlark.Thread,
			message string,
		) {
			fmt.Fprintln(a.stdout, message)
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

	predeclared := starlark.StringDict{
		"project": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(projectName),
				"path": starlark.String(projectPath),
			},
		),

		"template": starlarkstruct.FromStringDict(
			starlarkstruct.Default,
			starlark.StringDict{
				"name": starlark.String(template),
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

	options := &syntax.FileOptions{
		TopLevelControl: true,
	}

	_, err = starlark.ExecFileOptions(
		options,
		thread,
		hookPath,
		nil,
		predeclared,
	)

	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			return fmt.Errorf(
				"template hook failed:\n%s",
				evalErr.Backtrace(),
			)
		}

		return fmt.Errorf(
			"template hook failed: %w",
			err,
		)
	}

	return nil
}

func (h *hookContext) run(
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

	command, err := starlarkString(args[0])
	if err != nil {
		return nil, fmt.Errorf(
			"%s: program: %w",
			builtin.Name(),
			err,
		)
	}

	commandArgs := make(
		[]string,
		0,
		len(args)-1,
	)

	for i, value := range args[1:] {
		arg, err := starlarkString(value)

		if err != nil {
			return nil, fmt.Errorf(
				"%s: argument %d: %w",
				builtin.Name(),
				i+1,
				err,
			)
		}

		commandArgs = append(
			commandArgs,
			arg,
		)
	}

	cmd := exec.Command(
		command,
		commandArgs...,
	)

	cmd.Dir = h.projectPath
	cmd.Stdin = h.app.stdin
	cmd.Stdout = h.app.stdout
	cmd.Stderr = h.app.stderr

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

func (h *hookContext) replace(
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

	fullPath, err := h.resolvePath(path)
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

	text = strings.ReplaceAll(
		text,
		old,
		newValue,
	)

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(
		fullPath,
		[]byte(text),
		info.Mode().Perm(),
	); err != nil {
		return nil, fmt.Errorf(
			"%s: write %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.MakeInt(count), nil
}

func (h *hookContext) write(
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

	fullPath, err := h.resolvePath(path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(
		filepath.Dir(fullPath),
		0o755,
	); err != nil {
		return nil, fmt.Errorf(
			"%s: create parent directory: %w",
			builtin.Name(),
			err,
		)
	}

	if err := os.WriteFile(
		fullPath,
		[]byte(content),
		0o644,
	); err != nil {
		return nil, fmt.Errorf(
			"%s: write %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.None, nil
}

func (h *hookContext) mkdir(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var path string

	if err := starlark.UnpackArgs(
		builtin.Name(),
		args,
		kwargs,
		"path",
		&path,
	); err != nil {
		return nil, err
	}

	fullPath, err := h.resolvePath(path)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(
		fullPath,
		0o755,
	); err != nil {
		return nil, fmt.Errorf(
			"%s: %s: %w",
			builtin.Name(),
			path,
			err,
		)
	}

	return starlark.None, nil
}

func (h *hookContext) remove(
	_ *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var path string

	if err := starlark.UnpackArgs(
		builtin.Name(),
		args,
		kwargs,
		"path",
		&path,
	); err != nil {
		return nil, err
	}

	fullPath, err := h.resolvePath(path)
	if err != nil {
		return nil, err
	}

	if fullPath == h.projectPath {
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

func (h *hookContext) resolvePath(
	path string,
) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if filepath.IsAbs(path) {
		return "", fmt.Errorf(
			"path must be relative to project: %s",
			path,
		)
	}

	fullPath := filepath.Clean(
		filepath.Join(
			h.projectPath,
			path,
		),
	)

	relative, err := filepath.Rel(
		h.projectPath,
		fullPath,
	)

	if err != nil {
		return "", err
	}

	if relative == ".." ||
		strings.HasPrefix(
			relative,
			".."+string(filepath.Separator),
		) {

		return "", fmt.Errorf(
			"path escapes project directory: %s",
			path,
		)
	}

	return fullPath, nil
}

func starlarkString(
	value starlark.Value,
) (string, error) {
	stringValue, ok := starlark.AsString(value)

	if !ok {
		return "", fmt.Errorf(
			"expected string, got %s",
			value.Type(),
		)
	}

	return stringValue, nil
}
