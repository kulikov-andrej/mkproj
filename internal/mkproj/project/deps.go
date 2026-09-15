package project

import (
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/editor"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/setup"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/storage"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

var (
	getTemplate     = templates.Get
	createProject   = storage.Create
	runSetup        = setup.Run
	cleanupMetadata = storage.CleanupMetadata
	openEditor      = editor.Open
)
