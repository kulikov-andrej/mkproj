package core

import (
	"github.com/kulikov-andrej/mkproj/internal/mkproj/data/project"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/data/templates"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/editor"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/setup"
)

var (
	getTemplate     = templates.Get
	listTemplates   = templates.List
	createProject   = project.Create
	runSetup        = setup.Run
	openProject     = editor.Open
	cleanupMetadata = project.CleanupMetadata
)
