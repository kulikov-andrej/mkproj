package cli

import (
	"fmt"
	"io"
)

func showHelp(out io.Writer) {
	fmt.Fprintln(out, `Usage:
  mkproj -t <name> <path> [--open]
  mkproj --template=<name> <path> [--open]
  mkproj --list
  mkproj --help

Options:
  -t, --template <name>  Template to use
  -o, --open             Open project in Code
  -l, --list             List available templates
  -h, --help             Show help`)
}
