package main

import (
	"fmt"
	"os"
)

func main() {
	templateRoot, err := defaultTemplateRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app := newApp(
		templateRoot,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)

	if err := app.run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
