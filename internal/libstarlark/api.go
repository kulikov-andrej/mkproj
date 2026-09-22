package libstarlark

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type Function func(
	thread *starlark.Thread,
	builtin *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error)

type API struct {
	globals starlark.StringDict
}

func NewAPI() *API {
	return &API{
		globals: make(starlark.StringDict),
	}
}

func (api *API) Struct(
	name string,
	fields map[string]string,
) *API {
	values := make(starlark.StringDict, len(fields))

	for name, value := range fields {
		values[name] = starlark.String(value)
	}

	api.globals[name] = starlarkstruct.FromStringDict(
		starlarkstruct.Default,
		values,
	)

	return api
}

func (api *API) Function(
	name string,
	fn Function,
) *API {
	api.globals[name] = starlark.NewBuiltin(name, fn)
	return api
}

func (api *API) Value(
	name string,
	value starlark.Value,
) *API {
	api.globals[name] = value
	return api
}

func (api *API) Build() starlark.StringDict {
	globals := make(starlark.StringDict, len(api.globals))

	for name, value := range api.globals {
		globals[name] = value
	}

	return globals
}

func AsString(value starlark.Value) (string, error) {
	result, ok := starlark.AsString(value)
	if !ok {
		return "", fmt.Errorf(
			"expected string, got %s",
			value.Type(),
		)
	}

	return result, nil
}
