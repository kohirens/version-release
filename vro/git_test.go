package main

import (
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := git.CloneFromBundle(tt.repo, tmpDir, fixtureDir, ps)

			got, _ := IsChangelogUpToDate(repoDir, tt.chgLogFile)
			if got != tt.want {
				t.Errorf("IsChangelogUpToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
