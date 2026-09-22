package libstarlark

import (
	"fmt"

	"github.com/kulikov-andrej/mkproj/internal/libio"

	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

func ExecFile(
	path string,
	globals starlark.StringDict,
	streams libio.IO,
) (starlark.StringDict, error) {
	thread := &starlark.Thread{
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

	result, err := starlark.ExecFileOptions(
		&syntax.FileOptions{
			TopLevelControl: true,
		},
		thread,
		path,
		nil,
		globals,
	)
	if err != nil {
		return nil, formatError(err)
	}

	return result, nil
}

func formatError(err error) error {
	if evalErr, ok := err.(*starlark.EvalError); ok {
		return fmt.Errorf("%s", evalErr.Backtrace())
	}

	return err
}
