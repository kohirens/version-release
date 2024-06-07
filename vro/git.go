package main

import (
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/pkg/git"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"os"
	"strings"
)

// IsChangelogUpToDate Indicate if there are any changes to be added to the
// changelog.
func IsChangelogUpToDate(wd, chgLogFile string) (bool, error) {
	cliffConfigName := "cliff.toml"
	configFile := wd + "/" + cliffConfigName

	if !fsio.Exist(configFile) { // make a config when none present.
		if e := os.WriteFile(configFile, []byte(cliffConfig), 0776); e != nil {
			return true, e
		}
	}

	if e := gitcliff.BuildChangelog(wd, chgLogFile); e != nil {
		return true, e
	}

	// Check if the files for modification.
	status, e2 := git.StatusWithOptions(wd, []string{"--porcelain", "--", chgLogFile, cliffConfigName})
	if e2 != nil {
		return true, e2
	}

	log.Dbugf("git status output = %v", status)
	log.Logf("debug status = %v", status)

	return len(status) == 0, nil
}

func lastUpdateWasAutoChangelog(wd string) bool {
	lastLog := git.Log(wd, "HEAD")
	return strings.Contains(lastLog, "An automated update of CHANGELOG.md")
}

func NoChangesToRelease(wd string) bool {
	changes, e1 := gitcliff.UnreleasedChanges(wd)
	if e1 != nil {
		log.Errf(e1.Error())
		return false
	}

	return len(changes) == 0
}
