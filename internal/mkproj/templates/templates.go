package templates

import "github.com/kulikov-andrej/mkproj/internal/mkproj/templates/model"

type Template = model.Template

func Get(name string) (Template, error) {
	return getTemplate(name)
}

func List() ([]Template, error) {
	return listTemplates()
}
