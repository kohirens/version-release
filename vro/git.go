package main

import (
	"fmt"
	"github.com/kohirens/stdlib/cli"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/git"
	"github.com/kohirens/version-release-orb/vro/pkg/gitcliff"
	"github.com/kohirens/version-release-orb/vro/pkg/gittoolbelt"
	"strings"
)

func areChangelogChangesPresent(wd, chgLogFile string) (bool, error) {
	// Step 1: Run a command to update the changelog and see if the file shows
	// in the git status output.
	// NOTE: This will be an issue after publish changelog workflow runs
	// because it will produce changelog for a tag that has not been
	// released (we put the egg before the first chicken)
	isUpToDate, err3 := IsChangelogUpToDate(wd, chgLogFile)
	if err3 != nil {
		return false, err3
	}

	if !isUpToDate {
		// We now need to check if the changes are because of a recent "publish
		// changelog" merge that is intended for a future tag release; in
		// which case we want to correct behavior and skip another changelog
		// publish.

		sv, err4 := gittoolbelt.Semver(wd)
		if err4 != nil {
			return false, fmt.Errorf(stderr.SemverWithChgLogCheck, chgLogFile, err4.Error())
		}

		// This should produce the same changelog, even if it is for a future
		// release, and if it produces a new changelog, then ignore that.
		if e := gitcliff.RebuildChangelog(wd, chgLogFile, sv); e != nil {
			return false, fmt.Errorf(stderr.RebuildWithChgLogCheck, chgLogFile, e.Error())
		}

		status, err1 := git.Status(wd)
		if err1 != nil {
			return false, err1
		}

		// Check if the changes in the changelog are from the last commit.
		if !strings.Contains(string(status), chgLogFile) {
			return false, nil
		}
	}

	return !isUpToDate, nil
}

// IsChangelogUpToDate Indicate if there are any changes to be added to the
// changelog.
func IsChangelogUpToDate(wd, chgLogFile string) (bool, error) {
	// step 1: run command with no --bump,
	so, se, _, co := cli.RunCommand(
		wd,
		gitcliff.Cmd,
		[]string{"--output", chgLogFile},
	)

	log.Infof(stdout.Wd, wd)
	log.Infof(stdout.Cs, co)

	if se != nil {
		return true, fmt.Errorf(stderr.CouldNotUpdateChglog, so, se.Error())
	}

	// step 2: pull the git status, so we can check the state of the changelog
	status, err1 := git.Status(wd)
	if err1 != nil {
		return true, fmt.Errorf(stderr.GitStatus, status, err1.Error())
	}

	// step 3: check the git status output to see if the changelog is listed
	if status != nil && strings.Contains(string(status), chgLogFile) {
		log.Logf(stdout.FoundChgInFile, co, chgLogFile)
		// there are legit changes to add to the changelog.
		return false, nil
	}

	// step 4: nothing to do, tt truly is up-to-date and nothing needs to be done
	return true, nil
}

// IsCommitTagged Checks if there is a tag ref for a commit.
func IsCommitTagged(wd, commit string) bool {
	tag, foundTag := git.LookupTags(wd, commit)
	if foundTag {
		log.Logf(stdout.AlreadyReleased, commit, tag)
		return true
	}

	return false
}
