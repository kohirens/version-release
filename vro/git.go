package main

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/git"
	"github.com/kohirens/version-release-orb/vro/pkg/gitcliff"
)

// IsChangelogUpToDate Indicate if there are any changes to be added to the
// changelog. Side effect is that it will build the changelog.
func IsChangelogUpToDate(wd, chgLogFile string) (bool, error) {
	// step 1: run command with no --bump to build a new or update an existing log
	if e := gitcliff.BuildChangelog(wd, chgLogFile); e != nil {
		return true, e
	}

	// step 2: check the git status of changelog for modification
	status, e1 := git.StatusWithOptions(wd, []string{"--porcelain", chgLogFile})
	if e1 != nil {
		// we return true, but we don't really know since there was an error,
		// in any case don't do anything with the file, so pretend its up-to-date
		return true, fmt.Errorf(stderr.GitStatus, status, e1.Error())
	}

	log.Infof(stdout.GitStatus, string(status))

	git.PrintStatus(wd)

	// when git status --porcelain will be an empty string when the file is
	// update-to-date or the file does not exist.
	// however it should exist since we ran the Git-cliff command to build it.
	return len(status) == 0, nil
}
