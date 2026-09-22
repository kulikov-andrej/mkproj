package libio

import (
	"io"
)

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
