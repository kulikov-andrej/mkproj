package core

import (
	"github.com/kulikov-andrej/mkproj/internal/data/project"
	"github.com/kulikov-andrej/mkproj/internal/data/templates"
	"github.com/kulikov-andrej/mkproj/internal/editor"
	"github.com/kulikov-andrej/mkproj/internal/hooks"
)

var (
	getTemplate     = templates.Get
	listTemplates   = templates.List
	createProject   = project.Create
	runHook         = hooks.Run
	openProject     = editor.Open
	cleanupMetadata = project.CleanupMetadata
)
