package cli

import (
	"fmt"
	"strings"

	"github.com/kulikov-andrej/mkproj/internal/libio"
)

type projectArgs struct {
	template string
	target   string
	open     bool
}

func runProject(
	args []string,
	streams libio.IO,
) error {
	opts, err := parseProjectArgs(args)
	if err != nil {
		return err
	}

	if opts.template == "" {
		return fmt.Errorf("template is required")
	}

	target := opts.target
	if target == "" {
		target = "."
	}

	fmt.Fprintf(
		streams.Out,
		"Creating project %q using template %q...\n",
		target,
		opts.template,
	)

	proj, err := createProject(
		opts.template,
		target,
		streams,
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(
		streams.Out,
		"Created %q.\n  %s\n",
		proj.Name,
		proj.Path,
	)

	if opts.open {
		fmt.Fprintln(streams.Out, "Opening Code...")

		if err := openProject(proj); err != nil {
			return err
		}
	}

	fmt.Fprintln(streams.Out, "Have a vibey coding :)")

	return nil
}

func showProjectHelp(
	streams libio.IO,
) {
	fmt.Fprintln(streams.Out, `Usage:
  mkproj [<path>] -t <template> [--open]

Arguments:
  path                       Project path (default: current directory)

Options:
  -t, --template <template>  Template name to use
  -o, --open                 Open project in Code`)
}

func parseProjectArgs(
	args []string,
) (projectArgs, error) {
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
