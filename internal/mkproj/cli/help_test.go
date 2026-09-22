package cli

import (
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		contain string
	}{
		{name: "main", args: nil, contain: "Help topics:"},
		{name: "project", args: []string{"project"}, contain: "Project path (default: current directory)"},
		{name: "template", args: []string{"template"}, contain: "List available templates"},
		{name: "run", args: []string{"run"}, contain: "mkproj run <command>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			streams, out, _ := testStreams()
			if err := runHelp(tt.args, streams); err != nil {
				t.Fatalf("runHelp() error = %v", err)
			}
			if !strings.Contains(out.String(), tt.contain) {
				t.Fatalf("stdout = %q, want substring %q", out.String(), tt.contain)
			}
		})
	}
}

func TestRunHelpUnknownTopic(t *testing.T) {
	streams, _, _ := testStreams()

	err := runHelp([]string{"unknown"}, streams)
	if err == nil || err.Error() != `unknown topic "unknown"` {
		t.Fatalf("runHelp() error = %v", err)
	}
}
