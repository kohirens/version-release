package main

import (
	"github.com/kohirens/stdlib/test"
	"testing"
)

func TestIsChangelogUpToDate(t *testing.T) {
	tests := []struct {
		name       string
		chgLogFile string
		repo       string
		want       bool
	}{
		{"notUpToDate", "CHANGELOG.md", "repo-04", false},
		{"upToDate", "CHANGELOG.md", "repo-05", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := test.SetupARepository(tt.repo, tmpDir, fixtureDir, ps)

			got, _ := IsChangelogUpToDate(repoDir, tt.chgLogFile)
			if got != tt.want {
				t.Errorf("IsChangelogUpToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
