package main

import (
	"bufio"
	"fmt"
	"github.com/kohirens/stdlib"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/git"
	"github.com/kohirens/version-release-orb/vro/pkg/gitcliff"
	"os"
	"regexp"
	"strings"
)

// IsChangelogUpToDate Indicate if there are any changes to be added to the
// changelog.
func IsChangelogUpToDate(wd, chgLogFile string) (bool, error) {
	var version string
	// Check to see if the current changelog contains the unreleased tag.

	u, e1 := gitcliff.UnreleasedChangelogContext(wd)
	if e1 != nil {
		return false, e1
	}

	if len(u) > 0 {
		version = u[0].Version
		log.Logf(stdout.NextVersion, version)
	} else {
		version = git.GetCurrentTag(wd)
		log.Logf(stdout.CurrentVer, version)
	}

	if version == "" {
		return false, fmt.Errorf(stderr.NoVersionTag)
	}

	filename := wd + string(os.PathSeparator) + chgLogFile
	// there are unreleased changes and changelog does not exist
	if !stdlib.PathExist(filename) {
		return false, nil
	}

	chgLogRdr, e2 := os.Open(filename)
	if e2 != nil {
		return false, fmt.Errorf(stderr.OpenFile, chgLogFile, e2.Error())
	}

	defer chgLogRdr.Close()

	isUpToDate := false
	re := regexp.MustCompile(`^## \[v?` + version + `]\s+-\s+\d+-\d+-\d+`)
	scanner := bufio.NewScanner(chgLogRdr)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			log.Logf(stdout.Match)
			isUpToDate = true
			break
		}
	}

	return isUpToDate, nil
}

func lastUpdateWasAutoChangelog(wd string) bool {
	lasLog := git.LastLog(wd)
	return strings.Contains(lasLog, "An automated update of CHANGELOG.md")
}
