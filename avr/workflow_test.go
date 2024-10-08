package main

import (
	git2 "github.com/kohirens/stdlib/git"
	"github.com/kohirens/version-release/avr/pkg/git"
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
			repo := git2.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)
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

// Ensure the publish changelog works as advertised.
func TestPublishChangelog(t *testing.T) {
	ghcFixture, _ := newGitHubClient(&http.Client{})
	cases := []struct {
		name       string
		bundle     string
		chgLogFile string
		branch     string
		semVer     string
		ghc        GithubClient
		wantErr    bool
	}{
		{"successful", "repo-no-releases", "CHANGELOG.md", "main", "0.1.0", ghcFixture, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := git2.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)

			// set the default upstream for the fixture.
			oldUrl, _ := git.RemoteGetUrl(repo, "origin")
			_ = git.RemoteSetUrl(repo, "origin", "https://github.com/kohirens/repo-01", oldUrl)

			if err := PublishChangelog(repo, c.chgLogFile, c.branch, c.semVer, c.ghc); (err != nil) != c.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
