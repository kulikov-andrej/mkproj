package main

import (
	"path/filepath"
	"testing"
)

func TestTemplateHook(t *testing.T) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"cpp",
		map[string]string{
			"CMakeLists.txt": `
project({{PROJECT_NAME}})
`,
			".mkproj/hook.star": `
replace(
    "CMakeLists.txt",
    "{{PROJECT_NAME}}",
    project.name,
)

write(
    "generated.txt",
    project.name + "\n",
)
`,
		},
	)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	if err := env.app.run([]string{
		"-t",
		"cpp",
		target,
	}); err != nil {
		t.Fatal(err)
	}

	assertFileContent(
		t,
		filepath.Join(
			target,
			"CMakeLists.txt",
		),
		"\nproject(hello)\n",
	)

	assertFileContent(
		t,
		filepath.Join(
			target,
			"generated.txt",
		),
		"hello\n",
	)

	assertNotExists(
		t,
		filepath.Join(
			target,
			".mkproj",
		),
	)
}

func TestHookFailureStillRemovesMetadata(
	t *testing.T,
) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"broken",
		map[string]string{
			"hello.txt": "hello",
			".mkproj/hook.star": `
fail = missing_symbol
`,
		},
	)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	err := env.app.run([]string{
		"-t",
		"broken",
		target,
	})

	if err == nil {
		t.Fatal("expected hook failure")
	}

	assertNotExists(
		t,
		filepath.Join(
			target,
			".mkproj",
		),
	)
}

func TestHookCannotWriteOutsideProject(
	t *testing.T,
) {
	env := newTestApp(t)

	env.createTemplate(
		t,
		"evil",
		map[string]string{
			".mkproj/hook.star": `
write("../outside.txt", "oops")
`,
		},
	)

	target := filepath.Join(
		env.workRoot,
		"hello",
	)

	outside := filepath.Join(
		env.workRoot,
		"outside.txt",
	)

	err := env.app.run([]string{
		"-t",
		"evil",
		target,
	})

	if err == nil {
		t.Fatal("expected hook failure")
	}

	assertNotExists(t, outside)
}
