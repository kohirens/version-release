package main

import (
	"github.com/kohirens/stdlib/git"
	"net/http"
	"testing"
)

func TestIsChangelogCurrent(t *testing.T) {
	cases := []struct {
		name       string
		bundle     string
		chgLogFile string
		want       bool
		wantErr    bool
	}{
		{"noChangelog", "no-changelog", "CHANGELOG.md", false, true},
		{"notCurrent", "changelog-is-not-current", "CHANGELOG.md", false, false},
		{"current", "changelog-is-current", "CHANGELOG.md", true, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := git.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)
			got, _, err := isChangelogCurrent(repo, c.chgLogFile)

			if (err != nil) != c.wantErr {
				t.Errorf("IsChangelogCurrent() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if got != c.want {
				t.Errorf("IsChangelogCurrent() got = %v, want %v", got, c.want)
				return
			}
		})
	}
}

// Ensure the publishing a changelog works as advertised.
func TestPublishChangelog(runner *testing.T) {
	cases := []struct {
		name       string
		bundle     string
		chgLogFile string
		branch     string
		semVer     string
		wantErr    bool
	}{
		{"successful", "repo-no-releases", "CHANGELOG.md", "main", "0.1.0", false},
	}
	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			t.Setenv("CIRCLE_PROJECT_REPONAME", c.bundle)
			ghcFixture, _ := newGitHubClient(&http.Client{})
			repo := git.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)

			if err := PublishChangelog(repo, c.chgLogFile, c.branch, c.semVer, ghcFixture); (err != nil) != c.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}
