package scripts

import (
	"github.com/kulikov-andrej/mkproj/internal/libio"
)

type Context struct {
	ProjectName  string
	ProjectPath  string
	TemplateName string
	Streams      libio.IO
}
