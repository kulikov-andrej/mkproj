package core

import (
	"testing"

	"github.com/kulikov-andrej/mkproj/internal/data/templates"
)

func TestListTemplates(t *testing.T) {
	old := listTemplates

	expected := []templates.Template{
		{Name: "cpp"},
		{Name: "go"},
	}

	listTemplates = func() ([]templates.Template, error) {
		return expected, nil
	}

	t.Cleanup(func() {
		listTemplates = old
	})

	got, err := ListTemplates()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(expected) {
		t.Fatalf(
			"expected %#v, got %#v",
			expected,
			got,
		)
	}

	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf(
				"expected %#v, got %#v",
				expected,
				got,
			)
		}
	}
}
