package core

import (
	"github.com/kulikov-andrej/mkproj/internal/data/project"
	"github.com/kulikov-andrej/mkproj/internal/data/templates"
	"github.com/kulikov-andrej/mkproj/internal/editor"
	"github.com/kulikov-andrej/mkproj/internal/setup"
)

var (
	getTemplate     = templates.Get
	listTemplates   = templates.List
	createProject   = project.Create
	runSetup        = setup.Run
	openProject     = editor.Open
	cleanupMetadata = project.CleanupMetadata
)
