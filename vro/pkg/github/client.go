package github

import (
	"fmt"
	"github.com/kohirens/stdlib/log"
	"github.com/kohirens/version-release/vro/pkg/git"
	"io"
	"net/http"
	"time"
)

const (
	BaseUri         = "https://%s/repos/%s/%s"
	epBranches      = BaseUri + "/branches/%s"
	epPulls         = BaseUri + "/pulls"
	epPullMerge     = BaseUri + "/pulls/%d/merge"
	epRelease       = BaseUri + "/releases"
	epReleaseLatest = BaseUri + "/releases/latest"
	GenBranchName   = "auto-update-changelog"
	publicServer    = "github.com"
)

type Client struct {
	Client        HttpClient
	Domain        string
	MergeMethod   string
	Org           string
	RepositoryUri string
	Repository    string
	Token         string
	Username      string
	Host          string
}

func NewClient(repositoryUri, token, host string, client HttpClient) *Client {
	dom, org, repo := parseRepositoryUri(repositoryUri)

	if host == publicServer { // patch for public GitHub
		host = "api." + host
	} else {
		// patch for Enterprise server
		host = host + "/api/v3"
	}

	return &Client{
		Client:        client,
		Domain:        dom,
		Org:           org,
		Repository:    repo,
		RepositoryUri: repositoryUri,
		Token:         token,
		Username:      "git",
		Host:          host,
	}
}

func (gh *Client) DoesBranchExistRemotely(branch string) bool {
	uri := fmt.Sprintf(epBranches, gh.Host, gh.Org, gh.Repository, branch)

	res, err1 := gh.Send(uri, "GET", nil)
	if err1 != nil {
		log.Logf(stderr.CouldNotGetRequest, err1.Error())
		return false
	}

	bodyBits, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		log.Logf(stderr.CouldNotReadResponse, err2.Error())
		return false
	}

	fmt.Println(bodyBits)

	return res.StatusCode == 200
}

// PublishChangelog Stage, commit, and push local changes, then make a pull
// request and merge it containing:
// 1. the git-chglog command config files if they are missing or first run.
// 2. the CHANGELOG.md if it contains changes.
func (gh *Client) PublishChangelog(wd, branch, chaneLogFile string) error {
	if git.DoesBranchExistRemotely(wd, gh.RepositoryUri, GenBranchName) {
		return fmt.Errorf(
			stderr.BranchExists,
			GenBranchName, gh.RepositoryUri,
		)
	}

	if e := git.ConfigGlobal(wd, "user.name", gh.Username); e != nil {
		return e
	}

	email := fmt.Sprintf("%s@noreply.%s", gh.Username, gh.Domain)
	if e := git.ConfigGlobal(wd, "user.email", email); e != nil {
		return e
	}

	if e := git.CheckoutBranch(wd, GenBranchName); e != nil {
		return e
	}

	// Staging the CHANGELOG file.
	if e := git.StageFiles(wd, chaneLogFile); e != nil {
		return e
	}

	// Commit the CHANGELOG file.
	mergeBranchCommitMsg := "Updated the " + chaneLogFile
	mergeBranchCommitDesc := "An automated update of " + chaneLogFile
	if e := git.Commit(wd, mergeBranchCommitMsg, mergeBranchCommitDesc); e != nil {
		return e
	}

	if e := git.Push(wd, "origin", GenBranchName); e != nil {
		return e
	}

	pr, err1 := gh.OpenPullRequest(branch, GenBranchName, mergeBranchCommitMsg, mergeBranchCommitDesc)
	if err1 != nil {
		return err1
	}

	merge, err2 := gh.MergePullRequest(pr.Number, gh.MergeMethod)
	if err2 != nil {
		return err2
	}

	if merge.Merged {
		log.Logf(stdout.PullRequestMerged, pr.Number)
		return nil
	}

	log.Logf(stdout.MergeResponse, merge.Message)

	if e := gh.waitForPrToMerge(pr.Number, 5); e != nil {
		return e
	}

	return nil
}

func (gh *Client) Send(uri, method string, body io.Reader) (*http.Response, error) {
	req, err1 := http.NewRequest(method, uri, body)
	if err1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotPrepareRequest, err1.Error())
	}

	req.Header.Set("Authorization", "Bearer "+gh.Token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	res, err2 := gh.Client.Do(req)
	if err2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, method, err2.Error())
	}

	return res, nil
}

func (gh *Client) waitForPrToMerge(prNumber int, waitSeconds int) error {
	uri := fmt.Sprintf(epPullMerge, gh.Host, gh.Org, gh.Repository, prNumber)

	log.Logf(stdout.CheckMergeStatus, prNumber)

	res, err2 := gh.Send(uri, "GET", nil)
	if err2 != nil {
		return fmt.Errorf(stderr.CouldNotMakeRequest, "GET", err2.Error())
	}

	for i := 0; i < waitSeconds; i++ {
		time.Sleep(1 * time.Second)

		log.Infof("checking if pr %d was merged\n", prNumber)

		res, err2 = gh.Send(uri, "GET", nil)
		if err2 != nil {
			return fmt.Errorf(stderr.CouldNotPingMergeStatus, prNumber, err2.Error())
		}

		if res.StatusCode == 204 {
			break
		}
	}

	if res.StatusCode != 204 {
		return fmt.Errorf(stderr.MergeWaitTimeout, prNumber, waitSeconds)
	}

	return nil
}
