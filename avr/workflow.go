package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/kohirens/stdlib/fsio"
	"github.com/kohirens/version-release/avr/pkg/git"
	"github.com/kohirens/version-release/avr/pkg/gitcliff"
	"github.com/kohirens/version-release/avr/pkg/github"
	"os"
	"strings"
)

type GithubClient interface {
	TagAndRelease(revision, tag string) (*github.ReleasesResponse, error)
	PublishChangelog(wd, branch, header, msg, footer string, files []string) error
}

// PublishChangelog Run automation to update the CHANGELOG.md
func PublishChangelog(wd, chgLogFile, branch, semVer string, ghc GithubClient) error {
	files := []string{chgLogFile}

	configWritten, e1 := gitcliff.BuildConfig(wd + fsio.PS + gitcliff.CliffConfigName)
	if e1 != nil {
		return e1
	}

	if configWritten {
		files = append(files, gitcliff.CliffConfigName)
	}

	// If one exist, then checkout the changelog from the current tag and use that as the start point to prepend
	// unreleased changes to. This will keep GitCliff from pre-pending the same changes repeatedly to the changelog
	// in most cases. However, this will overwrite any manually added additions since the last tag, which should
	// be O.K. in most cases since this file is meant to be generated automatically.
	prevSemver := git.GetCurrentTag(wd)
	if prevSemver != "" {
		e2 := git.CheckoutFileFrom(wd, prevSemver, chgLogFile)
		if e2 != nil {
			return e2
		}
	}

	changes, hash, e3 := gitcliff.UnreleasedChangesHash(wd)
	if e3 != nil {
		log.Errf(e3.Error())
	}

	header := fmt.Sprintf(autoReleaseHeader, semVer)
	footer := fmt.Sprintf(autoReleaseFooter, base64.StdEncoding.EncodeToString(hash))

	// add the hash for the unreleased changes to the changelog for verification
	if e := os.Setenv("GIT_CLIFF__CHANGELOG__FOOTER", footer); e != nil {
		return fmt.Errorf(stderr.SetGitCliffFooter, e.Error())
	}

	// It is a guarantee that a changelog should always exist after running this build function.
	e4 := gitcliff.BuildChangelog(wd, chgLogFile, semVer)
	if e4 != nil {
		return e4
	}

	// Add additional files to the commit.
	addFiles := os.Getenv("PARAM_ADD_FILES_TO_COMMIT")
	if addFiles != "" {
		filesBytes, e := os.ReadFile(wd + "/" + addFiles)
		if e != nil {
			return fmt.Errorf(stderr.AdditionalFiles, e.Error())
		}
		cleanData := bytes.Replace(bytes.Trim(filesBytes, "\r\n"), []byte("\r\n"), []byte("\n"), -1)
		files = append(files, strings.Split(string(cleanData), "\n")...)
	}

	// Commit, push, and rebase the changelog.
	return ghc.PublishChangelog(wd, branch, header, string(changes), footer, files)
}

// changelogContains Return true if changelog contains the unreleased changes.
func changelogContains(wd, chgLogFile, subject string) (bool, error) {
	filename := wd + string(os.PathSeparator) + chgLogFile

	chgLogRdr, e2 := os.Open(filename)
	if e2 != nil {
		return false, fmt.Errorf(stderr.OpenFile, chgLogFile, e2.Error())
	}

	defer func() {
		e := chgLogRdr.Close()
		if e != nil {
			log.Errf(stderr.ClosingFile, filename, e.Error())
		}
	}()

	isUpToDate := false
	scanner := bufio.NewScanner(chgLogRdr)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, subject) {
			log.Infof(stdout.Match, subject)
			isUpToDate = true
			break
		}
	}

	return isUpToDate, nil
}

// isChangelogCurrent git-cliff blindly prepends commits to the CHANGELOG without
// verifying if they were already added. To prevent duplicate entries this function
// does special checks on the CHANGELOG
func isChangelogCurrent(wd, chgLogFile string) (bool, string, error) {
	// Get unreleased changes with the header and footer.
	changes, hash, e1 := gitcliff.UnreleasedChangesHash(wd)
	if e1 != nil {
		return false, "", e1
	}

	// When there are changes but the changelog is not present.
	if len(changes) > 0 && !fsio.Exist(wd+fsio.PS+chgLogFile) {
		return false, "", nil
	}

	log.Dbugf("len unreleased changes = %v", len(changes))

	hash64 := base64.StdEncoding.EncodeToString(hash)
	// Look for the hash footer in the last git log commit message, and also in the changelog, but before it is rebuilt.
	hashFooter := fmt.Sprintf(autoReleaseFooter, hash64)
	hasChanges, e3 := changelogContains(wd, chgLogFile, hashFooter)
	log.Dbugf("hashFooter: %v\n", hashFooter)
	log.Dbugf("hasChanges: %v\n", hasChanges)
	if e3 != nil {
		return false, "", e3
	}

	if !hasChanges {
		log.Infof(stderr.HashNotInChangelog, hashFooter)
		return false, "", nil
	}

	logMsg := git.LastLog(wd)
	log.Dbugf("logMsg: %v\n", logMsg)
	if !strings.Contains(logMsg, hashFooter) {
		log.Infof(stderr.HashNotInLog, hashFooter)
		return false, "", nil
	}

	return true, hash64, nil
}
