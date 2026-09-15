package cli

import (
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

var (
	listTemplates = templates.List
	createProject = project.Create
	openProject   = project.Open
)
