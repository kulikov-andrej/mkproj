package main

import (
	"fmt"
	"io"
	"strings"
)

type options struct {
	template string
	target   string
	open     bool
	list     bool
	help     bool
}

func parseArgs(args []string) (options, error) {
	var opts options

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if after, ok := strings.CutPrefix(arg, "--template="); ok {
			opts.template = after

			if strings.TrimSpace(opts.template) == "" {
				return opts, fmt.Errorf("--template requires a name")
			}

			continue
		}

		switch arg {
		case "-t", "--template":
			if i+1 >= len(args) {
				return opts, fmt.Errorf("%s requires a template name", arg)
			}

			i++
			opts.template = args[i]

			if strings.TrimSpace(opts.template) == "" {
				return opts, fmt.Errorf("%s requires a template name", arg)
			}

		case "-o", "--open":
			opts.open = true

		case "-l", "--list":
			opts.list = true

		case "-h", "--help":
			opts.help = true

		default:
			if strings.HasPrefix(arg, "-") {
				return opts, fmt.Errorf("unknown option: %s", arg)
			}

			if opts.target != "" {
				return opts, fmt.Errorf("unexpected argument: %s", arg)
			}

			opts.target = arg
		}
	}

	return opts, nil
}

func showHelp(w io.Writer, templateRoot string) {
	fmt.Fprintf(w, `Usage:
  mkproj -t <name> <path> [--open]
  mkproj --template=<name> <path> [--open]

Options:
  -t, --template <name>   Project template
  -o, --open              Open project in VS Code
  -l, --list              List available templates
  -h, --help              Show this help

Examples:
  mkproj --list
  mkproj --template=python backend --open
  mkproj --template=python .

Templates:
  %s
`, templateRoot)
}
