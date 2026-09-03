package main

import (
	"fmt"
	"io"
	"strings"
)

type app struct {
	templateRoot string

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	openProject func(string) error
}

func newApp(
	templateRoot string,
	stdin io.Reader,
	stdout io.Writer,
	stderr io.Writer,
) *app {
	app := &app{
		templateRoot: templateRoot,
		stdin:        stdin,
		stdout:       stdout,
		stderr:       stderr,
	}

	app.openProject = app.openInCode

	return app
}

func (a *app) run(args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	if len(args) == 0 || opts.help {
		showHelp(a.stdout, a.templateRoot)
		return nil
	}

	if opts.list {
		return a.showTemplates()
	}

	if strings.TrimSpace(opts.template) == "" {
		return fmt.Errorf(
			"missing template; use -t <name> or --template=<name>",
		)
	}

	if strings.TrimSpace(opts.target) == "" {
		return fmt.Errorf("missing project path")
	}

	return a.newProject(
		opts.template,
		opts.target,
		opts.open,
	)
}
