package scripts

import "io"

type Context struct {
	ProjectName  string
	ProjectPath  string
	TemplateName string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}
