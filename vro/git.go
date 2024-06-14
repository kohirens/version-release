package main

import (
	"bufio"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/pkg/git"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"os"
	"regexp"
)

// IsChangelogUpToDate Indicate if there are any changes to be added to the
// changelog.
func IsChangelogUpToDate(wd, chgLogFile, branch string) (bool, error) {
	hasUnreleasedChanges, e4 := gitcliff.UnreleasedChanges(wd)
	if e4 != nil {
		return true, e4
	}

	if hasUnreleasedChanges == nil || len(hasUnreleasedChanges) == 0 {
		return true, nil
	}
	// Scan the changelog to verify it does not already contain the unreleased entries.
	// git-cliff just blindly prepends commit to the CHANGELOG without verify they were already added. So we want to prevent duplicate entries.
	if fsio.Exist(chgLogFile) {
		// Exit when the change is update-to-date or an error occurred
		if containUnreleased, e5 := changelogContains(&hasUnreleasedChanges[0], wd, chgLogFile); containUnreleased || e5 != nil {
			return true, e5
		}
	}

	return false, nil
}

// changelogContains Return true if changelog contains the unreleased changes.
func changelogContains(unreleased *gitcliff.Unreleased, wd, chgLogFile string) (bool, error) {
	filename := wd + string(os.PathSeparator) + chgLogFile

	chgLogRdr, e2 := os.Open(filename)
	if e2 != nil {
		return false, fmt.Errorf(stderr.OpenFile, chgLogFile, e2.Error())
	}

	defer chgLogRdr.Close()

	isUpToDate := false
	re := regexp.MustCompile(`^## \[v?` + unreleased.Version + `]\s+-\s+\d+-\d+-\d+`)
	scanner := bufio.NewScanner(chgLogRdr)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			log.Logf(stdout.Match, unreleased.Version)
			isUpToDate = true
			break
		}
	}

	return isUpToDate, nil
}

// Old_IsChangelogUpToDate Indicate if there are any changes to be added to the
// changelog.
func Old_IsChangelogUpToDate(wd, chgLogFile string) (bool, error) {
	_, e := gitcliff.BuildChangelog(wd, chgLogFile)
	if e != nil {
		return true, e
	}

	// Check if the files for modification.
	status, e2 := git.StatusWithOptions(wd, []string{"--porcelain", "--", chgLogFile, gitcliff.CliffConfigName})
	if e2 != nil {
		return true, e2
	}

	log.Dbugf(stdout.GitStatus, status)

	utd := len(status) == 0

	if utd {
		log.Logf(stdout.ChgLogUpToDate)
	} else {
		log.Logf(stdout.ChgLogNotUpToDate, status)
	}

	return utd, nil
}
