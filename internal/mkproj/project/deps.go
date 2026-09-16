package project

import (
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/editor"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/scripts"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/project/storage"
	"github.com/kulikov-andrej/mkproj/internal/mkproj/templates"
)

var (
	getTemplate       = templates.Get
	createProject     = storage.Create
	getCurrentProject = storage.GetCurrentProject
	loadScript        = scripts.Load
	callScript        = scripts.Call
	cleanupMetadata   = storage.CleanupMetadata
	openEditor        = editor.Open
)
