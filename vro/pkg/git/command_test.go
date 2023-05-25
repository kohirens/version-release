package git

import (
	help "github.com/kohirens/stdlib/test"
	"os"
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
			repo := help.SetupARepository(tt.bundle, tmpDir, fixtureDir, ps)
			if got := DoesBranchExistRemotely(".", repo, tt.branch); got != tt.want {
				t.Errorf("DoesBranchExistRemotely() = %v, want %v", got, tt.want)
			}
		})
	}
}
