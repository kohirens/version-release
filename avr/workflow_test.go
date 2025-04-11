package main

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
	"github.com/kohirens/version-release/avr/pkg/github"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"net/http"
	"os"
	"testing"
	"time"
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
			help.Chdir(t, repo)
			_ = os.WriteFile(repo+ps+"CHANGELOG.md", []byte("[1.0.0] - 2024-06-14\n\n### Added\n\n- README.md"), 0664)

			if err := PublishChangelog(repo, c.chgLogFile, c.branch, c.semVer, ghcFixture); (err != nil) != c.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}

// Ensure publishing a changelog can also take a list of files, and adds them
// to the commit of the updated CHANGELOG.md.
func TestPublishChangelogWithAdditionalFiles(runner *testing.T) {
	cases := []struct {
		name       string
		bundle     string
		chgLogFile string
		branch     string
		semVer     string
		files      map[string][]byte
		addFiles   map[string][]byte
		wantErr    bool
	}{
		{
			"successful",
			"readme-only-no-tags",
			"CHANGELOG.md",
			"main",
			"",
			map[string][]byte{
				"CHANGELOG.md":       []byte("[0.1.0] - 2025-04-11\n\n### Added\n\n- README.md\n- generated.txt\n"),
				"generated-file.txt": []byte("generated-file.txt\n\n- README.md"),
			},
			map[string][]byte{
				"files-to-add.txt": []byte("generated-file.txt\n"),
			},
			false,
		},
	}
	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			ghcFixture := &github.Client{
				Client: &http.Client{
					Timeout: 5 * time.Second,
				},
				Domain:      "github.com",
				MergeMethod: "rebase",
				Repository:  "kohirens/" + c.bundle,
				Token:       "fakegithubtoken",
				Host:        "https://api.github.com",
			}
			repo := git.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)

			_ = os.WriteFile(repo+ps+"CHANGELOG.md", []byte("[1.0.0] - 2024-06-14\n\n### Added\n\n- README.md"), 0664)
			lib.MakeFiles(repo, c.files)
			lib.MakeFiles(repo, c.addFiles)

			t.Setenv("PARAM_ADD_FILES_TO_COMMIT", "files-to-add.txt")
			if err := PublishChangelog(repo, c.chgLogFile, c.branch, c.semVer, ghcFixture); (err != nil) != c.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}
