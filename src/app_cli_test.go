package main

import (
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	env := newTestApp(t)

	if err := env.app.run([]string{"--help"}); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(
		env.stdout.String(),
		"Usage:",
	) {
		t.Fatalf(
			"expected help output, got %q",
			env.stdout.String(),
		)
	}
}

func TestNoArgumentsShowsHelp(t *testing.T) {
	env := newTestApp(t)

	if err := env.app.run(nil); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(
		env.stdout.String(),
		"Usage:",
	) {
		t.Fatalf(
			"expected help output, got %q",
			env.stdout.String(),
		)
	}
}

func TestUnknownOption(t *testing.T) {
	env := newTestApp(t)

	err := env.app.run([]string{"--wat"})

	if err == nil {
		t.Fatal("expected an error")
	}

	if !strings.Contains(
		err.Error(),
		"unknown option",
	) {
		t.Fatalf("unexpected error: %v", err)
	}
}
