package gitcliff

import (
	"bytes"
	"fmt"
	"github.com/kohirens/stdlib"
	help "github.com/kohirens/stdlib/test"
	"os"
	"testing"
)

const (
	fixtureDir = "testdata"
	tmpDir     = "tmp"
	ps         = string(os.PathSeparator)
)

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

			if bytes.Compare(actual, expected) != 0 {
				t.Errorf("BuildChangelog() error %v does not match expected output", chgLogFile)
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
