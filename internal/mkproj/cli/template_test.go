package cli

import (
	"errors"
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func TestRunTemplateWithoutCommandShowsHelp(t *testing.T) {
	streams, out, _ := testStreams()

	if err := runTemplate(nil, streams); err != nil {
		t.Fatalf("runTemplate() error = %v", err)
	}
	if got := out.String(); got == "" {
		t.Fatal("stdout is empty, want template help")
	}
}

func TestRunTemplateList(t *testing.T) {
	old := listTemplates
	listTemplates = func() ([]templates.Template, error) {
		return []templates.Template{{Name: "go"}, {Name: "rust"}}, nil
	}
	t.Cleanup(func() { listTemplates = old })

	streams, out, errOut := testStreams()
	if err := runTemplate([]string{"list"}, streams); err != nil {
		t.Fatalf("runTemplate() error = %v", err)
	}

	if got, want := out.String(), "go\nrust\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if errOut.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", errOut.String())
	}
}

func TestRunTemplateListEmpty(t *testing.T) {
	old := listTemplates
	listTemplates = func() ([]templates.Template, error) { return nil, nil }
	t.Cleanup(func() { listTemplates = old })

	streams, out, errOut := testStreams()
	if err := runTemplate([]string{"list"}, streams); err != nil {
		t.Fatalf("runTemplate() error = %v", err)
	}

	if out.Len() != 0 {
		t.Fatalf("stdout = %q, want empty", out.String())
	}
	if got, want := errOut.String(), "No templates found.\n"; got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}

func TestRunTemplateListReturnsError(t *testing.T) {
	wantErr := errors.New("list failed")
	old := listTemplates
	listTemplates = func() ([]templates.Template, error) { return nil, wantErr }
	t.Cleanup(func() { listTemplates = old })

	streams, _, _ := testStreams()
	err := runTemplate([]string{"list"}, streams)
	if !errors.Is(err, wantErr) {
		t.Fatalf("runTemplate() error = %v, want %v", err, wantErr)
	}
}

func TestRunTemplateRejectsInvalidCommand(t *testing.T) {
	streams, _, _ := testStreams()

	tests := []struct {
		args    []string
		wantErr string
	}{
		{args: []string{"unknown"}, wantErr: `unknown template command "unknown"`},
		{args: []string{"list", "extra"}, wantErr: `unexpected argument "extra"`},
	}

	for _, tt := range tests {
		err := runTemplate(tt.args, streams)
		if err == nil || err.Error() != tt.wantErr {
			t.Fatalf("runTemplate(%v) error = %v, want %q", tt.args, err, tt.wantErr)
		}
	}
}
