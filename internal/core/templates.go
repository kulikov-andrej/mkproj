package core

import (
	"github.com/kulikov-andrej/mkproj/internal/data/templates"
)

func ListTemplates() ([]templates.Template, error) {
	return listTemplates()
}
