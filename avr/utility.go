package main

import (
	"bufio"
	"fmt"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/avr/pkg/git"
	"github.com/kohirens/version-release/avr/pkg/gitcliff"
	"github.com/kohirens/version-release/avr/pkg/github"
	"github.com/kohirens/version-release/avr/pkg/lib"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// changelogContainsUnreleased Return true if changelog contains the unreleased changes.
func changelogContainsUnreleased(unreleased *gitcliff.Unreleased, wd, chgLogFile string) (bool, error) {
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

// TagIt Only consider tagging if HEAD has no tag and the commit message
// contains the expected auto-release header.
func TagIt(wd, commit, semVer string) bool {
	hasSemverTag := git.HasSemverTag(wd, commit)

	// Log that the commit already has a tag.
	if hasSemverTag {
		log.Logf(stderr.CommitAlreadyTagged, commit)
		return false
	}

	nextVer := gitcliff.NextVersion(wd, semVer)
	if nextVer == "" { // No version to tag, then check for changelog updates.
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

func newGitHubClient(client *http.Client) (*github.Client, error) {
	var gh *github.Client

	switch clo.CiCd {
	case circleci.Name:
		eVars, e1 := lib.GetRequiredEnvVars([]string{
			github.EnvToken,
			circleci.EnvProjectRepoName,
			circleci.EnvProjectUsername,
			github.EnvApiUrl,
		})

		if e1 != nil {
			return nil, e1
		}

		gh = github.NewClient(
			eVars[circleci.EnvProjectUsername]+"/"+eVars[circleci.EnvProjectRepoName],
			eVars[github.EnvToken],
			eVars[github.EnvApiUrl],
			client,
		)
	case github.Name:
		eVars, e1 := lib.GetRequiredEnvVars([]string{
			github.EnvApiUrl,
			github.EnvRepository,
			github.EnvRepositoryOwner,
			github.EnvToken,
		})

		if e1 != nil {
			return nil, e1
		}

		gh = github.NewClient(
			eVars[github.EnvRepository],
			eVars[github.EnvToken],
			eVars[github.EnvApiUrl],
			client,
		)

	}

	return gh, nil
}
