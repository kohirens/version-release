package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

/*
	curl -L \
	  -H "Accept: application/vnd.github+json" \
	  -H "Authorization: Bearer <YOUR-TOKEN>" \
	  -H "X-GitHub-Api-Version: 2022-11-28" \
	  https://api.github.com/repos/OWNER/REPO/pulls
*/
const ( // See https://docs.github.com/en/rest?apiVersion=2022-11-28 for endpoints and permissions.
	BaseUri        = "%s/repos/%s"
	epBranches     = BaseUri + "/branches/%s"
	epCommit       = BaseUri + "/git/commits"
	epCreateARef   = BaseUri + "/git/refs"
	epGetARef      = BaseUri + "/git/ref/%v"
	epUpdateARef   = BaseUri + "/git/refs/%v"
	epPulls        = BaseUri + "/pulls"
	epPullMerge    = BaseUri + "/pulls/%d/merge"
	epRelease      = BaseUri + "/releases"
	epReleaseByTag = epRelease + "/tags/%s"
	epTree         = BaseUri + "/git/trees"
	fullRefPrefix  = "refs/heads/"
	GenBranchName  = "auto-update-changelog"
	repoUrl        = "git@%v:%v.git"
)

type Client struct {
	Client      HttpClient
	Domain      string
	MergeMethod string
	Repository  string
	Server      string
	Token       string
	Host        string
}

// NewClient GitHub API client. Set the host with the `-gh-api-url` options.
// Public host: https://api.github.com
// Possible enterprise format: https://<github-enterprise-server>/api/v3
func NewClient(repo, token, host string, client HttpClient) *Client {

	return &Client{
		Client:     client,
		Repository: repo,
		Token:      token,
		Host:       host,
	}
}

// DoesBranchExistRemotely Check the GitHub API that a branch exist remotely.
func (gh *Client) DoesBranchExistRemotely(branch string) bool {
	uri := fmt.Sprintf(epBranches, gh.Host, gh.Repository, branch)

	res, err1 := gh.Send(uri, "GET", nil)
	if err1 != nil {
		Log.Logf(stderr.CouldNotGetRequest, err1.Error())
		return false
	}

	bodyBits, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		Log.Logf(stderr.CouldNotReadResponse, err2.Error())
		return false
	}

	Log.Dbugf(stdout.RemoteBranchStatus, string(bodyBits))

	return res.StatusCode == 200
}

// PublishChangelog Stage, commit, and push local changes, then make a pull
// request and merge it containing the CHANGELOG.md if it contains changes.
func (gh *Client) PublishChangelog(wd, baseRef, header, msgBody, footer string) error {
	// Return early if the branch that updates the change log exists remotely.
	uri := fmt.Sprintf(repoUrl, PublicServer, gh.Repository)
	Log.Dbugf(stdout.RepoUrl, uri)

	if gh.ReferenceExist("heads/" + GenBranchName) {
		return fmt.Errorf(stderr.BranchExists, GenBranchName, uri)
	}

	if e := Push(wd, GenBranchName, baseRef, header+"\n"+msgBody+"\n"+footer, gh); e != nil {
		return e
	}

	pr, e1 := gh.OpenPullRequest(baseRef, GenBranchName, header, msgBody+"\n"+footer)
	if e1 != nil {
		return e1
	}
	merge, e2 := gh.MergePullRequest(pr.Number, gh.MergeMethod)
	if e2 != nil {
		return e2
	}

	if merge.Merged {
		Log.Logf(stdout.PullRequestMerged, pr.Number)
		return nil
	}

	Log.Logf(stdout.MergeResponse, merge.Message)

	if e := gh.waitForPrToMerge(pr.Number, 5); e != nil {
		return e
	}

	return nil
}

// ReleaseByTag Check the GitHub API that a branch exist remotely.
func (gh *Client) ReleaseByTag(tag string) (*ReleasesResponse, error) {
	uri := fmt.Sprintf(epReleaseByTag, gh.Host, gh.Repository, tag)

	Log.Dbugf(stdout.RequestUrl, uri)

	res, e1 := gh.Send(uri, "GET", nil)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotGetRequest, e1.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(
			stderr.ResponseStatusCode,
			res.StatusCode,
			uri,
			res.Status,
		)
	}

	bodyBytes, e2 := io.ReadAll(res.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e2.Error())
	}

	rr := &ReleasesResponse{}
	e3 := json.Unmarshal(bodyBytes, rr)
	if e3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e3.Error())
	}

	return rr, nil
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
	uri := fmt.Sprintf(epPullMerge, gh.Host, gh.Repository, prNumber)

	Log.Logf(stdout.CheckMergeStatus, uri)

	r1, e1 := gh.Send(uri, "GET", nil)
	if e1 != nil {
		return fmt.Errorf(stderr.CouldNotMakeRequest, "GET", e1.Error())
	}

	if r1.StatusCode != 204 {
		return fmt.Errorf(stderr.MergeWaitTimeout, prNumber, waitSeconds)
	}

	var statusCode int

	for i := 0; i < waitSeconds; i++ {
		time.Sleep(1 * time.Second)

		Log.Infof("checking if pr %d was merged\n", prNumber)
		r2, e2 := gh.Send(uri, "GET", nil)
		if e2 != nil {
			return fmt.Errorf(stderr.CouldNotPingMergeStatus, prNumber, e2.Error())
		}

		statusCode = r2.StatusCode
		if statusCode == 204 {
			break
		}
	}

	if statusCode != 204 {
		return fmt.Errorf(stderr.MergeWaitTimeout, prNumber, waitSeconds)
	}

	return nil
}
