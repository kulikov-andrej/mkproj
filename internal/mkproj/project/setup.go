package project

import (
	"os"
	"path/filepath"

	"github.com/kulikov-andrej/mkproj/internal/libio"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/scripts"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

func runTemplateSetup(
	proj Project,
	template templates.Template,
	streams libio.IO,
) error {
	setupPath := filepath.Join(
		proj.Path,
		".mkproj",
		"setup.star",
	)

	if _, err := os.Stat(setupPath); os.IsNotExist(err) {
		return nil
	}

	ctx := scripts.Context{
		ProjectName:  proj.Name,
		ProjectPath:  proj.Path,
		TemplateName: template.Name,
		Streams:      streams,
	}

	return scripts.ExecFile(setupPath, ctx)
}
