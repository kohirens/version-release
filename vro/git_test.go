package main

import (
	"github.com/kohirens/stdlib/git"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"testing"
)

func TestIsChangelogUpToDate(t *testing.T) {
	tests := []struct {
		name       string
		chgLogFile string
		repo       string
		semVer     string
		want       bool
	}{
		{"has-changes", "CHANGELOG.md", "repo-04", "", false},
		{"no-changes", "CHANGELOG.md", "repo-05", "", true},
		{"repo-has-chglog-changes", "CHANGELOG.md", "repo-has-chglog-changes", "", false},
		{"repo-has-no-chglog-changes", "CHANGELOG.md", "repo-has-no-chglog-changes", "", true},
		{"repo-no-previous-tag", "CHANGELOG.md", "repo-no-previous-tag", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoDir := git.CloneFromBundle(tt.repo, tmpDir, fixtureDir, ps)

			got, _ := IsChangelogUpToDate(repoDir, tt.chgLogFile, tt.semVer)
			if got != tt.want {
				t.Errorf("IsChangelogUpToDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_changelogContains Verify repos contains or does not contain unreleased
// changes.
func Test_changelogContains(t *testing.T) {
	tests := []struct {
		name       string
		unreleased *gitcliff.Unreleased
		repo       string
		want       bool
		wantErr    bool
	}{
		{"contains-unreleased", &gitcliff.Unreleased{Version: "0.1.0"}, "repo-chglog-with-unreleased", true, false},
		{"does-not-contains", &gitcliff.Unreleased{Version: "0.1.1"}, "repo-chglog-without-unreleased", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wd := git.CloneFromBundle(tt.repo, tmpDir, fixtureDir, ps)
			got, err := changelogContains(tt.unreleased, wd, "CHANGELOG.md")
			if (err != nil) != tt.wantErr {
				t.Errorf("changelogContains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("changelogContains() got = %v, want %v", got, tt.want)
			}
		})
	}
}
