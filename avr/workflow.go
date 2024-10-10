package main

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/avr/pkg/circleci"
	"github.com/kohirens/version-release/vro/pkg/gitcliff"
	"net/http"
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

func circleCiTagWorkflow(branch string, client *http.Client) error {
	ciClient := circleci.NewClient("gh", clo.Project, clo.Username, client)

	// Build pipeline parameters to trigger the tag-and-release workflow.
	pp, e3 := circleci.PipelineParameters(branch, publishReleaseTagWorkflow)
	if e3 != nil {
		return e3

	}

	log.Logf(stdout.TriggerWorkflow, publishReleaseTagWorkflow)

	//  Trigger the workflow
	return ciClient.TriggerWorkflow(pp)
}
