package github

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
	vgit "github.com/kohirens/version-release/vro/pkg/git"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	ps = string(os.PathSeparator)
)

var tmpDir = help.AbsPath("tmp")

func TestMain(m *testing.M) {
	help.ResetDir(tmpDir, 0777)
	os.Exit(m.Run())
}

func TestClient_PublishChangelog(t *testing.T) {
	tests := []struct {
		name    string
		bundle  string
		branch  string
		header  string
		msgBody string
		files   []string
		wantErr bool
	}{
		{"success", "repo-commit-message", "main", "auto: Release 1.0.0", "## Added README.md", []string{"CHANGELOG.md"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//gh := github.NewClient("https://github.com/kohirens/repo-01", "fakegithubtoken", "github.com", mockClient)
			gh := &Client{
				Client: &http.Client{
					Timeout: 5 * time.Second,
				},
				Domain:        "github.com",
				MergeMethod:   "rebase",
				Org:           "kohirens",
				RepositoryUri: "https://github.com/kohirens/repo-01",
				Repository:    "repo-01",
				Token:         "fakegithubtoken",
				Username:      "RobotTest",
				Host:          "api.github.com",
			}
			repo := git.CloneFromBundle(tt.bundle, "tmp", "testdata", ps)
			oldUrl, _ := vgit.RemoteGetUrl(repo, "origin")
			_ = vgit.RemoteSetUrl(repo, "origin", "https://github.com/kohirens/repo-01", oldUrl)
			_ = os.WriteFile(repo+ps+"CHANGELOG.md", []byte("[1.0.0] - 2024-06-14\n\n### Added\n\n- README.md"), 0664)
			_ = vgit.StageFiles(repo, "CHANGELOG.md")

			if err := gh.PublishChangelog(repo, tt.branch, tt.header, tt.msgBody, tt.files); (err != nil) != tt.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
