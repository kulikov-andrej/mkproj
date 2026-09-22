package templates

import (
	"os"
)

var (
	resolveRoot = getDefaultRoot
	readDir     = os.ReadDir
)
