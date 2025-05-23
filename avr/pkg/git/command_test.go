package git

import (
	"github.com/kohirens/stdlib/git"
	help "github.com/kohirens/stdlib/test"
	"os"
	"reflect"
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

func TestCommit(t *testing.T) {
	tests := []struct {
		name    string
		bundle  string
		message string
		wantErr bool
	}{
		{"multiline-commit-message", "repo-commit-message", "auto: Release 0.1.0\n\n## Added\n* README.md", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := git.CloneFromBundle(tt.bundle, tmpDir, fixtureDir, ps)

			_ = Config(repo, "user.name", "Robot Test")
			_ = Config(repo, "user.email", "test@example.com")
			_ = os.WriteFile(repo+ps+"testfile.txt", []byte("Salam"), 0664)
			_ = StageFiles(repo, "testfile.txt")

			if err := Commit(repo, tt.message); (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}

			got := Log(repo, "HEAD")
			if !strings.Contains(got, "auto: Release 0.1.0") {
				t.Errorf("Commit()\n%v\n\tdoe NOT contain: %q", got, tt.message)
			}
		})
	}
}

func TestCheckoutFileFrom(t *testing.T) {
	cases := []struct {
		name     string
		bundle   string
		treeIsh  string
		filename string
		wantErr  bool
	}{
		{"noPreviousTag", "no-previous-tag", "0.1.0", "README.md", true},
		{"hasPreviousTag", "has-1-tag", "0.1.0", "README.md", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := git.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)

			err := CheckoutFileFrom(repo, c.treeIsh, c.filename)

			if (err != nil) != c.wantErr {
				t.Errorf("CheckoutFileFrom() error = %v, wantErr %v", err, c.wantErr)
				return
			}
		})
	}
}

func TestGetCurrentTag(t *testing.T) {
	cases := []struct {
		name   string
		bundle string
		want   string
	}{
		{"hasNoTag", "no-previous-tag", ""},
		{"hasTag", "has-1-tag", "0.1.0"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			repo := git.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)
			if got := GetCurrentTag(repo); got != c.want {
				t.Errorf("GetCurrentTag() = %v, want %v", got, c.want)
			}
		})
	}
}

func TestStatus(runner *testing.T) {
	runner.Skipf("skip for now, not moving forward with using Status to detect files to add to the changelog commit")
	cases := []struct {
		name   string
		bundle string
		files  map[string]string
		want   *StatusPorcelainFiles
	}{
		{
			"add-new-files",
			"git-status-porcelain-readme-only",
			map[string]string{"CHANGELOG.md": "[0.1.0] - 2025-04-12"},
			&StatusPorcelainFiles{},
		},
	}

	for _, c := range cases {
		runner.Run(c.name, func(t *testing.T) {
			repo := git.CloneFromBundle(c.bundle, tmpDir, fixtureDir, ps)
			var fil []*File

			fil = append(fil, &File{
				Stage: "?",
				Path:  "CHANGELOG.md",
			})
			makeFiles(repo, c.files)
			if got := Status(repo); !reflect.DeepEqual(got, c.want) {
				t.Errorf("Status() = %v, want %v", got, c.want)
			}
		})
	}
}

func makeFiles(wd string, files map[string]string) {
	for filename, content := range files {
		if e := os.WriteFile(wd+ps+filename, []byte(content), 0664); e != nil {
			panic(e)
		}
	}
}
