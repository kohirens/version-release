package gittoolbelt

import (
	help "github.com/kohirens/stdlib/test"
	"os"
)

const (
	ps = string(os.PathSeparator)
)

var (
	tmpDir     = help.AbsPath("tmp")
	fixtureDir = "testdata"
)
