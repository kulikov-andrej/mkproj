package scripts

import "github.com/kulikov-andrej/mkproj/internal/libstarlark"

func ExecFile(
	path string,
	ctx Context,
) error {
	api := libstarlark.NewAPI().
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

	_, err := libstarlark.ExecFile(path, api, ctx.Streams)
	return err
}
