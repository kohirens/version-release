package main

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release-orb/vro/pkg/circleci"
	"github.com/kohirens/version-release-orb/vro/pkg/gittoolbelt"
)

type Workflow struct {
	// GitHubClient GitHub API client
	GitHubClient circleci.GithubClient
	// Token A CircleCI API token
	Token string
}

func NewWorkflow(token string, ghClient circleci.GithubClient) *Workflow {
	return &Workflow{
		GitHubClient: ghClient,
		Token:        token,
	}
}

// PublishChangelog Run automation to update the CHANGELOG.md
func (wf *Workflow) PublishChangelog(wd, chgLogFile, branch string) error {
	// Step 1: Determine if the changelog has updates
	isUpdated, err1 := IsChangelogUpToDate(wd, chgLogFile)
	if err1 != nil {
		return err1
	}

	// nothing to publish
	if isUpdated {
		// If there were no changes to publish then how did we get here?
		// This pipeline should not have been triggered.
		return fmt.Errorf(stderr.NoChangelogChanges)
	}

	// Step 4: Commit, push, and publish the changelog.
	return wf.GitHubClient.PublishChangelog(wd, branch, chgLogFile)
}

// PublishReleaseTag Publish a release on GitHub.
func (wf *Workflow) PublishReleaseTag(branch, wd string) error {
	// Step 1: Grab semantic version info.
	si, err1 := gittoolbelt.Semver(wd)
	if err1 != nil {
		return fmt.Errorf(stderr.CouldNotGetVersion, err1)
	}

	// Step 2: Publish a new tag on GitHub.
	rr, err2 := wf.GitHubClient.TagAndRelease(branch, si)
	if err2 != nil {
		return err2
	}

	log.Logf(stdout.ReleaseTag, rr.Name)

	return nil
}
