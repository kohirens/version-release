package main

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
)

type Workflow struct {
	// GitHubClient GitHub API client
	GitHubClient circleci.GithubClient
	// Token A CircleCI API token
}

func NewWorkflow(ghClient circleci.GithubClient) *Workflow {
	return &Workflow{
		GitHubClient: ghClient,
	}
}

// PublishChangelog Run automation to update the CHANGELOG.md
func (wf *Workflow) PublishChangelog(wd, chgLogFile, branch, semVer string) error {
	files, e1 := gitcliff.BuildChangelog(wd, chgLogFile, semVer)
	if e1 != nil {
		return e1
	}

	msg, e2 := gitcliff.UnreleasedMessage(wd)
	if e2 != nil {
		log.Errf(e2.Error())
	}

	header := fmt.Sprintf(autoReleaseHeader, semVer)
	// Commit, push, and rebase the changelog.
	return wf.GitHubClient.PublishChangelog(wd, branch, header, string(msg), files)
}

// PublishReleaseTag Publish a release on GitHub.
func (wf *Workflow) PublishReleaseTag(branch, semVer string) error {

	return nil
}
