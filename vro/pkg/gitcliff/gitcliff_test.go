package gitcliff

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib"
	help "github.com/kohirens/stdlib/test"
	"os"
	"testing"
	"time"
)

const (
	fixtureDir = "testdata"
	tmpDir     = "tmp"
	ps         = string(os.PathSeparator)
)

func TestMain(m *testing.M) {
	// Only runs when this environment variable is set.
	help.ResetDir(tmpDir, 0777)

	// Run all tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestBuildChangelog(t *testing.T) {
	chgLogFile := "CHANGELOG.md"
	tests := []struct {
		name    string
		bundle  string
		wantErr bool
	}{
		{"no-changelog", "repo-01", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := help.SetupARepository(tt.bundle, tmpDir, fixtureDir, ps)

			if err := BuildChangelog(repo, chgLogFile); (err != nil) != tt.wantErr {
				t.Errorf("BuildChangelog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			actual := loadFile(repo + ps + chgLogFile)
			expected := loadFile(fixtureDir + ps + fmt.Sprintf("%v-expected-chglog.txt", tt.bundle))
			expt := fmt.Sprintf(string(expected), time.Now().Format("2006-01-02"))

			if bytes.Compare(actual, []byte(expt)) != 0 {
				t.Errorf("BuildChangelog() error %v does not match expected output", chgLogFile)
				return
			}
		})
	}
}

func TestHasUnreleasedChanges(t *testing.T) {
	tests := []struct {
		name    string
		bundle  string
		want    bool
		wantErr bool
	}{
		{"has-unreleased", "repo-02", true, false},
		{"no-unreleased", "repo-03", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := help.SetupARepository(tt.bundle, tmpDir, fixtureDir, ps)
			got, e1 := HasUnreleasedChanges(repo)
			if (e1 != nil) != tt.wantErr {
				t.Errorf("HasUnreleasedChanges() error = %v, wantErr %v", e1.Error(), tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("HasUnreleasedChanges() want %v, got %v", got, tt.want)
				return
			}
		})
	}
}

func loadFile(filename string) []byte {
	if !stdlib.PathExist(filename) {
		panic(fmt.Sprintf("file %v not found", filename))
	}
	contents, _ := os.ReadFile(filename)
	return contents
}

func TestBump(t *testing.T) {
	tests := []struct {
		name   string
		bundle string
		want   string
	}{
		{"has-unreleased-commits", "repo-04", "0.1.0"},
		{"no-unreleased", "repo-05", ""},
		{"major-release", "repo-06", "1.0.0"},
		{"minor-release", "repo-07", "0.2.0"},
		{"patch-release", "repo-08", "0.1.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := help.SetupARepository(tt.bundle, tmpDir, fixtureDir, ps)

			if got := Bump(repo); got != tt.want {
				t.Errorf("Bump() = %v, want %v", got, tt.want)
			}
		})
	}
}
