package gittoolbelt

import (
	"github.com/kohirens/stdlib/git"
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

func TestIsTaggable(t *testing.T) {
	tests := []struct {
		name   string
		bundle string
		want   bool
	}{
		{"is-taggable", "repo-02", true},
		{"not-taggable", "repo-03", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)
			got := IsTaggable(repo)

			if got != tt.want {
				t.Errorf("IsTaggable() want %v, got %v", got, tt.want)
				return
			}
		})
	}
}
