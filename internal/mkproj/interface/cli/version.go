package cli

import (
	"fmt"
	"io"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/buildinfo"
)

func runVersion(out io.Writer) error {
	fmt.Fprintf(
		out,
		"mkproj %s\n",
		buildinfo.Version,
	)

	return nil
}
