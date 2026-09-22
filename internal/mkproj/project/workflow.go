package project

import (
	"path/filepath"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/scripts"
)

func WorkflowCommands(
	proj Project,
	streams libio.IO,
) ([]string, error) {
	return scripts.Commands(
		workflowPath(proj),
		workflowContext(proj, streams),
	)
}

func RunWorkflow(
	proj Project,
	command string,
	streams libio.IO,
) error {
	return scripts.Run(
		workflowPath(proj),
		command,
		workflowContext(proj, streams),
	)
}

func workflowPath(
	proj Project,
) string {
	return filepath.Join(
		proj.Path,
		".mkproj",
		"workflow.star",
	)
}

func workflowContext(proj Project, streams libio.IO) scripts.Context {
	return scripts.Context{
		ProjectName: proj.Name,
		ProjectPath: proj.Path,
		Streams:     streams,
	}
}
