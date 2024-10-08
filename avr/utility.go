package main

import (
	"bufio"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/pkg/git"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"os"
	"regexp"
	"strings"
)

const CheckSemVer = `v?\d+\.\d+\.\d+(-.+)?`

// changelogContains Return true if changelog contains the unreleased changes.
func changelogContains(unreleased *gitcliff.Unreleased, wd, chgLogFile string) (bool, error) {
	filename := wd + string(os.PathSeparator) + chgLogFile

	chgLogRdr, e2 := os.Open(filename)
	if e2 != nil {
		return false, fmt.Errorf(stderr.OpenFile, chgLogFile, e2.Error())
	}

	defer func() {
		e := chgLogRdr.Close()
		if e != nil {
			log.Errf("Error closing chgLog file: %s", e.Error())
		}
	}()

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

func validateMergeType(mType string) (string, error) {
	if mType == "" {
		return "", fmt.Errorf(stderr.MergeTypeEmpty)
	}

	if mType != "rebase" && mType != "merge" && mType != "squash" {
		return "", fmt.Errorf(stderr.MergeType)
	}

	return mType, nil
}

func nextVersion(semVer string, wd string) (string, error) {
	// check if a version has been provided as input.
	nextVer := semVer
	if nextVer == "" {
		nextVer = gitcliff.Bump(wd)
	}

	log.Infof(stdout.NextSemVer, nextVer)

	if nextVer == "" {
		return "", fmt.Errorf(stderr.NothingToTag)
	}

	return nextVer, nil
}

// TagIt Only consider tagging if HEAD has no tag and the commit message
// contains the expected auto-release header.
func TagIt(wd, commit, semVer string) bool {
	hasSemverTag := git.HasSemverTag(wd, commit)

	// Log that the commit already has a tag.
	if hasSemverTag {
		log.Logf(stderr.CommitAlreadyTagged, commit)
		return false
	}

	nextVer, e2 := nextVersion(semVer, wd)
	if e2 != nil { // No version to tag, then check for changelog updates.
		log.Logf(e2.Error())
		return false
	}

	l := git.Log(wd, commit)
	log.Dbugf(stdout.DbgCommitLog, l)

	// Skip commits that are NOT released and also should NOT be tagged.
	if !strings.Contains(l, fmt.Sprintf(autoReleaseHeader, nextVer)) {
		return true
	}

	return true
}
