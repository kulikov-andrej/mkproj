package cli

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/buildinfo"
)

func TestRunVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		[]string{"version"},
		libio.IO{In: &bytes.Buffer{}, Out: &stdout, Err: &stderr},
	)
	if err != nil {
		t.Fatal(err)
	}

	want := fmt.Sprintf(
		"mkproj %s\n",
		buildinfo.Version,
	)

	if got := stdout.String(); got != want {
		t.Fatalf(
			"expected stdout %q, got %q",
			want,
			got,
		)
	}

	if stderr.Len() != 0 {
		t.Fatalf(
			"unexpected stderr: %q",
			stderr.String(),
		)
	}
}
