package main

import (
	"fmt"
	"os"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/cli"
)

func main() {
	if err := cli.Run(
		os.Args[1:],
		libio.IO{
			In:  os.Stdin,
			Out: os.Stdout,
			Err: os.Stderr,
		},
	); err != nil {
		fmt.Fprintln(os.Stderr, "mkproj:", err)
		os.Exit(1)
	}
}
