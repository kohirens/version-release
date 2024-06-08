package main

import (
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/git"
	"testing"
)

func TestIsChangelogUpToDate(t *testing.T) {
	tests := []struct {
		name       string
		chgLogFile string
		repo       string
		want       bool
	}{
		{"has-changes", "CHANGELOG.md", "repo-04", false},
		{"no-changes", "CHANGELOG.md", "repo-05", true},
		{"repo-has-chglog-changes", "CHANGELOG.md", "repo-has-chglog-changes", false},
		{"repo-has-no-chglog-changes", "CHANGELOG.md", "repo-has-no-chglog-changes", true},
		{"repo-no-previous-tag", "CHANGELOG.md", "repo-no-previous-tag", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := git.CloneFromBundle(tt.repo, tmpDir, fixtureDir, ps)

			got, _ := IsChangelogUpToDate(repoDir, tt.chgLogFile)
			if got != tt.want {
				t.Errorf("IsChangelogUpToDate() = %v, want %v", got, tt.want)
			}

			if !fsio.Exist(repoDir + "/cliff.toml") {
				t.Errorf("cliff.toml does not exist")
			}
		})
	}
}
