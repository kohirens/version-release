package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Verification struct {
	Verified   bool   `json:"verified"`
	Reason     string `json:"reason"`
	Signature  string `json:"signature"`
	Payload    string `json:"payload"`
	VerifiedAt string `json:"verified_at"`
}

type UserInfo struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Commit struct {
	Author       UserInfo     `json:"author"`
	Committer    UserInfo     `json:"committer"`
	Message      string       `json:"message"`
	Tree         Tree         `json:"tree"`
	Url          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `json:"verification"`
}

type Profile struct {
	Login             string `json:"login"`
	Id                int    `json:"id"`
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string `json:"url"`
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsUrl  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type ParentCommit struct {
	Sha     string `json:"sha"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
}

type Link struct {
	Self string `json:"self"`
	Html string `json:"html"`
}

type RequiredStatusChecks struct {
	EnforcementLevel string        `json:"enforcement_level"`
	Contexts         []interface{} `json:"contexts"`
	Checks           []interface{} `json:"checks"`
}

type Protection struct {
	Enabled              bool                 `json:"enabled"`
	RequiredStatusChecks RequiredStatusChecks `json:"required_status_checks"`
}

type CommitDetail struct {
	Sha         string         `json:"sha"`
	NodeId      string         `json:"node_id"`
	Commit      Commit         `json:"commit"`
	Url         string         `json:"url"`
	HtmlUrl     string         `json:"html_url"`
	CommentsUrl string         `json:"comments_url"`
	Author      Profile        `json:"author"`
	Committer   Profile        `json:"committer"`
	Parents     []ParentCommit `json:"parents"`
}

type ResponseBranch struct {
	Name          string       `json:"name"`
	Commit        CommitDetail `json:"commit"`
	Links         Link         `json:"_links"`
	Protected     bool         `json:"protected"`
	Protection    Protection   `json:"protection"`
	ProtectionUrl string       `json:"protection_url"`
}

// GitReference Git references within a repository.
type GitReference struct {
	Ref    string `json:"ref"`
	NodeId string `json:"node_id"`
	Url    string `json:"url"`
	Object struct {
		Type string `json:"type"`
		Sha  string `json:"sha"`
		Url  string `json:"url"`
	} `json:"object"`
}

type RequestGitReference struct {
	Ref string `json:"ref"`
	Sha string `json:"sha"`
}

type RequestUpdateReference struct {
	Sha   string `json:"sha"`
	Force bool   `json:"force"`
}

// ReferenceExist Verify a reference exist in GitHub.
func (gh *Client) ReferenceExist(ref string) bool {
	_, e1 := GetReference(ref, gh)
	if e1 != nil {
		return false
	}
	return true
}

// GetReference Make a new branch using the GitHub API
// `get /repos/{owner}/{repo}/git/ref/{ref}` endpoint.
// Returns a single reference from your Git database. The :ref in the URL must
// be formatted as `heads/<branch name>` for branches and `tags/<tag name>` for
// tags. If the :ref doesn't match an existing ref, a 404 is returned.
// For details see https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28#get-a-reference
func GetReference(ref string, gh *Client) (*GitReference, error) {
	uri := fmt.Sprintf(epGetARef, gh.Host, gh.Repository, ref)

	Log.Dbugf("get reference uri: %v", uri)
	Log.Dbugf("ref: %v", ref)

	res, e1 := gh.Send(uri, "GET", nil)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, e1.Error())
	}

	resBody, e2 := io.ReadAll(res.Body)
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e2.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(stderr.GetRef, res.Status, string(resBody))
	}

	gr := &GitReference{}
	if e4 := json.Unmarshal(resBody, gr); e4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e4.Error())
	}

	Log.Dbugf("gr: %v", gr)

	return gr, nil
}

// NewReference Make a new reference using the GitHub API
// `post /repos/{owner}/{repo}/git/refs` endpoint.
// You are unable to create new references for empty repositories, even if the
// commit SHA-1 hash used exists. Empty repositories are repositories without
// branches.
// For details see
// https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28#create-a-reference
func NewReference(name, sha string, gh *Client) (*GitReference, error) {
	uri := fmt.Sprintf(epCreateARef, gh.Host, gh.Repository)

	Log.Dbugf("new reference uri: %v", uri)
	Log.Dbugf("name: %v", name)
	Log.Dbugf("sha: %v", sha)

	body := &RequestGitReference{name, sha}

	reqBody, e1 := json.Marshal(body)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, e1.Error())
	}

	res, e2 := gh.Send(uri, "POST", bytes.NewReader(reqBody))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, e2.Error())
	}

	resBody, e3 := io.ReadAll(res.Body)
	if e3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e3.Error())
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf(stderr.NewRef, res.Status, string(resBody))
	}

	gr := &GitReference{}
	if e4 := json.Unmarshal(resBody, gr); e4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e4.Error())
	}

	return gr, nil
}

// UpdateReference Make a new branch using the GitHub API
// `patch /repos/{owner}/{repo}/git/refs/{ref}` endpoint.
// For details see
// https://docs.github.com/en/rest/git/refs?apiVersion=2022-11-28#update-a-reference
//
//	Note there are some inconsistencies in the GitHub REST API, not sure about
//	its GraphQL. You should use the same ref format as GetReference, or you may
//	get a 404. Not sure why they did it this way, it cost some time to figure
//	this out on your own.
func UpdateReference(ref, sha string, force bool, gh *Client) (*GitReference, error) {
	uri := fmt.Sprintf(epUpdateARef, gh.Host, gh.Repository, ref)

	Log.Dbugf("update reference uri: %v", uri)
	Log.Dbugf("ref: %v", ref)
	Log.Dbugf("sha: %v", sha)
	Log.Dbugf("force: %v", force)

	body := &RequestUpdateReference{sha, force}

	reqBody, e1 := json.Marshal(body)
	if e1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, e1.Error())
	}

	res, e2 := gh.Send(uri, "PATCH", bytes.NewReader(reqBody))
	if e2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, e2.Error())
	}

	resBody, e3 := io.ReadAll(res.Body)
	if e3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, e3.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(stderr.UpdateRef, res.Status, string(resBody))
	}

	gr := &GitReference{}
	if e4 := json.Unmarshal(resBody, gr); e4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, e4.Error())
	}

	return gr, nil
}
