package github

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
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
		files   []string
		wantErr bool
	}{
		{
			"success",
			"github-repo-commit-message",
			"main",
			"auto: Release 1.0.0",
			"## Added README.md",
			[]string{"CHANGELOG.md"},
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
			help.Chdir(t, repo)
			// we need a CHANGELOG.md fixture.
			_ = os.WriteFile("CHANGELOG.md", []byte("[1.0.0] - 2024-06-14\n\n### Added\n\n- README.md"), 0664)

			if err := gh.PublishChangelog(c.branch, c.header, c.msgBody, "", c.files); (err != nil) != c.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, c.wantErr)
			}
		})
	}
}
