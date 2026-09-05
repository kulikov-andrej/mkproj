package main

import (
	"fmt"
	"os"

	"github.com/kulikov-andrej/mkproj/internal/interface/cli"
)

func main() {
	if err := cli.Run(
		os.Args[1:],
		os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintln(os.Stderr, "mkproj:", err)
		os.Exit(1)
	}
}
