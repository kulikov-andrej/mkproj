package main

import (
	"fmt"
	"os"

	"github.com/kulikov-andrej/mkproj/internal/installer"
)

func main() {
	if err := installer.Run(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "mkproj-install:", err)
		os.Exit(1)
	}
}
