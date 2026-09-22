package project

import (
	"strings"
	"testing"
)

func TestOpenReportsMissingCode(t *testing.T) {
	t.Setenv("PATH", "")

	err := Open(Project{Path: t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "code executable not found") {
		t.Fatalf("Open() error = %v, want code executable not found", err)
	}
}
