package libstarlark

import (
	"testing"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestAPI(t *testing.T) {
	fn := func(
		*starlark.Thread,
		*starlark.Builtin,
		starlark.Tuple,
		[]starlark.Tuple,
	) (starlark.Value, error) {
		return starlark.None, nil
	}

	globals := NewAPI().
		Struct("project", map[string]string{
			"name": "hello",
			"path": "/tmp/hello",
		}).
		Function("run", fn).
		Value("enabled", starlark.True).
		Build()

	project, ok := globals["project"].(*starlarkstruct.Struct)
	if !ok {
		t.Fatalf(
			"project = %T, want *starlarkstruct.Struct",
			globals["project"],
		)
	}

	name, err := project.Attr("name")
	if err != nil {
		t.Fatal(err)
	}

	if got, _ := starlark.AsString(name); got != "hello" {
		t.Fatalf("project.name = %q, want %q", got, "hello")
	}

	run, ok := globals["run"].(*starlark.Builtin)
	if !ok {
		t.Fatalf("run = %T, want *starlark.Builtin", globals["run"])
	}

	if got := run.Name(); got != "run" {
		t.Fatalf("run.Name() = %q, want %q", got, "run")
	}

	if globals["enabled"] != starlark.True {
		t.Fatalf("enabled = %v, want True", globals["enabled"])
	}
}

func TestAsStringRejectsNonString(t *testing.T) {
	_, err := AsString(starlark.MakeInt(1))
	if err == nil {
		t.Fatal("AsString() error = nil, want error")
	}
}
