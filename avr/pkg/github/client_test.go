package github

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
	vgit "github.com/kohirens/version-release/avr/pkg/git"
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
		{"success", "github-repo-commit-message", "main", "auto: Release 1.0.0", "## Added README.md", []string{"CHANGELOG.md"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gh := &Client{
				Client: &http.Client{
					Timeout: 5 * time.Second,
				},
				Domain:      "github.com",
				MergeMethod: "rebase",
				Repository:  "kohirens/" + tt.bundle,
				Token:       "fakegithubtoken",
				Host:        "https://api.github.com",
			}
			repo := git.CloneFromBundle(tt.bundle, "tmp", "testdata", ps)
			_ = os.WriteFile(repo+ps+"CHANGELOG.md", []byte("[1.0.0] - 2024-06-14\n\n### Added\n\n- README.md"), 0664)
			_ = vgit.StageFiles(repo, "CHANGELOG.md")

			if err := gh.PublishChangelog(tt.branch, tt.header, tt.msgBody, "", tt.files); (err != nil) != tt.wantErr {
				t.Errorf("PublishChangelog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
