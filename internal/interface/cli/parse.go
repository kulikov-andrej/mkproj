package cli

import (
	"fmt"
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

		switch {
		case arg == "-h" || arg == "--help":
			opts.help = true

		case arg == "-l" || arg == "--list":
			opts.list = true

		case arg == "-o" || arg == "--open":
			opts.open = true

		case arg == "-t" || arg == "--template":
			if i+1 >= len(args) {
				return options{}, fmt.Errorf(
					"%s requires a value",
					arg,
				)
			}

			i++
			opts.template = args[i]

			if strings.TrimSpace(opts.template) == "" {
				return options{}, fmt.Errorf(
					"%s requires a value",
					arg,
				)
			}

		case strings.HasPrefix(arg, "--template="):
			opts.template = strings.TrimPrefix(
				arg,
				"--template=",
			)

			if strings.TrimSpace(opts.template) == "" {
				return options{}, fmt.Errorf(
					"--template requires a value",
				)
			}

		case strings.HasPrefix(arg, "-"):
			return options{}, fmt.Errorf(
				"unknown option: %s",
				arg,
			)

		default:
			if strings.TrimSpace(opts.target) != "" {
				return options{}, fmt.Errorf(
					"unexpected argument: %s",
					arg,
				)
			}

			opts.target = arg
		}
	}

	return opts, nil
}
