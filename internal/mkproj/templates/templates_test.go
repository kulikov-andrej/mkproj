package templates

import "testing"

func resetDependencies(t *testing.T) {
	t.Helper()

	oldGetTemplate := getTemplate
	oldListTemplates := listTemplates

	t.Cleanup(func() {
		getTemplate = oldGetTemplate
		listTemplates = oldListTemplates
	})
}

func TestGet(t *testing.T) {
	resetDependencies(t)

	want := Template{
		Name: "cpp",
		Path: "templates/cpp",
	}

	getTemplate = func(name string) (Template, error) {
		if name != want.Name {
			t.Fatalf(
				"expected name %q, got %q",
				want.Name,
				name,
			)
		}

		return want, nil
	}

	got, err := Get(want.Name)
	if err != nil {
		t.Fatal(err)
	}

	if got != want {
		t.Fatalf(
			"expected %#v, got %#v",
			want,
			got,
		)
	}
}

func TestList(t *testing.T) {
	resetDependencies(t)

	want := []Template{
		{Name: "cpp"},
		{Name: "go"},
	}

	listTemplates = func() ([]Template, error) {
		return want, nil
	}

	got, err := List()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != len(want) {
		t.Fatalf(
			"expected %#v, got %#v",
			want,
			got,
		)
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf(
				"expected %#v, got %#v",
				want,
				got,
			)
		}
	}
}
