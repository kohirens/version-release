package github

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	ps = string(os.PathSeparator)
)

var tmpDir = "tmp"

func TestMain(m *testing.M) {
	help.ResetDir(tmpDir, 0777)
	os.Exit(m.Run())
}

func TestClient_PublishChangelog(runner *testing.T) {
	cases := []struct {
		name    string
		bundle  string
		branch  string
		header  string
		msgBody string
		files   map[string][]byte
		wantErr bool
	}{
		{
			"success",
			"github-repo-commit-message",
			"main",
			"auto: Release 1.0.0",
			"## Added README.md",
			map[string][]byte{"CHANGELOG.md": []byte("[1.0.0] - 2024-06-14\n\n### Added\n\n- README.md")},
			false,
		},
	}
	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			gh := &Client{
				Client: &http.Client{
					Timeout: 5 * time.Second,
				},
				Domain:      "github.com",
				MergeMethod: "rebase",
				Repository:  "kohirens/" + c.bundle,
				Token:       "fakegithubtoken",
				Host:        "https://api.github.com",
			}
			repo := git.CloneFromBundle(c.bundle, "tmp", "testdata", ps)
			fileFixtures := lib.MakeFiles(repo, c.files)
			if err := gh.PublishChangelog(repo, c.branch, c.header, c.msgBody, "", fileFixtures); (err != nil) != c.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
