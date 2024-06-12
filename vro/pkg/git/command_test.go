package git

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
	"os"
	"strings"
	"testing"
)

const (
	ps = string(os.PathSeparator)
)

var tmpDir = help.AbsPath("tmp")
var fixtureDir = "testdata"

func TestMain(m *testing.M) {
	help.ResetDir(tmpDir, 0777)
	os.Exit(m.Run())
}

func TestDoesBranchExistRemotely(t *testing.T) {
	tests := []struct {
		name   string
		branch string
		bundle string
		want   bool
	}{
		{"exist", "main", "repo-02", true},
		{"notExist", "main2", "repo-02", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)
			if got := DoesBranchExistRemotely(".", repo, tt.branch); got != tt.want {
				t.Errorf("DoesBranchExistRemotely() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasSemverTag(t *testing.T) {
	tests := []struct {
		name   string
		bundle string
		commit string
		want   bool
	}{
		{"no-semver", "repo-08", "11f23fc5a62476ba57def51d1d7e8c020926d26c", false},
		{"has-semver", "repo-09", "82edbde9750818d6312cf18ea11f1be8525d5e47", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)

			if got := HasSemverTag(repo, tt.commit); got != tt.want {
				t.Errorf("hasSemverTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLog(t *testing.T) {
	tests := []struct {
		name   string
		bundle string
		refID  string
		want   string
	}{
		{"success", "repo-03", "HEAD", "test1234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)

			if got := Log(repo, tt.refID); !strings.Contains(got, tt.want) {
				t.Errorf("LastLog() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCommit(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name   string
		bundle string
		commit string
		want   bool
	}{
		{"is-commit", "repo-03", "HEAD", true},
		{"not-a-commit", "repo-03", "abcd123", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)

			if got := IsCommit(repo, tt.commit); got != tt.want {
				t.Errorf("IsCommit() = %v, want %v", got, tt.want)
			}
		})
	}
}
