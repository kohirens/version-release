package gittoolbelt

import (
	help "github.com/kohirens/stdlib/test"
	"os"
	"testing"
)

const (
	ps = string(os.PathSeparator)
)

var (
	tmpDir     = help.AbsPath("tmp")
	fixtureDir = "testdata"
)

func TestAddGitChgLogConfig(t *testing.T) {
	tests := []struct {
		name       string
		bundle     string
		configFile string
		want       bool
	}{
		{"added", "repo-01", ".chglog/config.yml", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := help.SetupARepository(tt.bundle, tmpDir, fixtureDir, ps)
			if got := AddGitChgLogConfig(wd, tt.configFile, ""); got != tt.want {
				t.Errorf("AddGitChgLogConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
