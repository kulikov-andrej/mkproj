package cli

import (
	"fmt"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/buildinfo"
)

func runVersion(
	streams libio.IO,
) error {
	fmt.Fprintf(streams.Out, "mkproj %s\n", buildinfo.Version)

	return nil
}
