package cli

import (
	"fmt"
	"strings"
)

type projectArgs struct {
	template string
	target   string
	open     bool
}

func parseProjectArgs(args []string) (projectArgs, error) {
	var options projectArgs

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-o" || arg == "--open":
			options.open = true

		case arg == "-t" || arg == "--template":
			if i+1 >= len(args) {
				return projectArgs{}, fmt.Errorf(
					"%s requires a value",
					arg,
				)
			}

			i++
			options.template = args[i]

			if strings.TrimSpace(options.template) == "" {
				return projectArgs{}, fmt.Errorf(
					"%s requires a value",
					arg,
				)
			}

		case strings.HasPrefix(arg, "--template="):
			options.template = strings.TrimPrefix(
				arg,
				"--template=",
			)

			if strings.TrimSpace(options.template) == "" {
				return projectArgs{}, fmt.Errorf(
					"--template requires a value",
				)
			}

		case strings.HasPrefix(arg, "-"):
			return projectArgs{}, fmt.Errorf(
				"unknown option: %s",
				arg,
			)

		default:
			if strings.TrimSpace(options.target) != "" {
				return projectArgs{}, fmt.Errorf(
					"unexpected argument: %s",
					arg,
				)
			}

			options.target = arg
		}
	}

	return options, nil
}
