package core

import (
	"github.com/kulikov-andrej/mkproj/internal/mkproj/data/templates"
)

func ListTemplates() ([]templates.Template, error) {
	return listTemplates()
}
