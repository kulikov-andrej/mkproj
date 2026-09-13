package cli

import (
	"fmt"
	"io"
)

func showHelp(out io.Writer) {
	fmt.Fprintln(out, `Usage:
  mkproj [<path>] -t <template> [--open]
  mkproj [<path>] --template=<template> [--open]
  mkproj --list
  mkproj --help
  mkproj --version

Options:
  -t, --template <template>  Template name to use
  -o, --open                 Open project in Code
  -l, --list                 List available templates
  -h, --help                 Show help
  -v, --version              Show current version`)
}
