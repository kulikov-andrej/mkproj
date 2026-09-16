package cli

import (
	"bytes"
	"strings"
	"testing"
)

type runResult struct {
	stdout string
	stderr string
	err    error
}

func resetDependencies(t *testing.T) {
	t.Helper()

	oldListTemplates := listTemplates
	oldCreateProject := createProject
	oldOpenProject := openProject
	oldGetCurrentProject := getCurrentProject
	oldListWorkflowCommands := listWorkflowCommands
	oldRunProjectWorkflow := runProjectWorkflow

	t.Cleanup(func() {
		listTemplates = oldListTemplates
		createProject = oldCreateProject
		openProject = oldOpenProject
		getCurrentProject = oldGetCurrentProject
		listWorkflowCommands = oldListWorkflowCommands
		runProjectWorkflow = oldRunProjectWorkflow
	})
}

func runCLI(args ...string) runResult {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	err := Run(
		args,
		strings.NewReader(""),
		&stdout,
		&stderr,
	)

	return runResult{
		stdout: stdout.String(),
		stderr: stderr.String(),
		err:    err,
	}
}
