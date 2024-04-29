package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"io"
	"time"
)

type Release struct {
	Body         string `json:"body"`
	Name         string `json:"name"`
	TagName      string `json:"tag_name"`
	TargetCommit string `json:"target_commitish"`
}

type Uploader struct {
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

type Arthor struct {
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

type Asset struct {
	Url                string    `json:"url"`
	BrowserDownloadUrl string    `json:"browser_download_url"`
	Id                 int       `json:"id"`
	NodeId             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	State              string    `json:"state"`
	ContentType        string    `json:"content_type"`
	Size               int       `json:"size"`
	DownloadCount      int       `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Uploader           Uploader  `json:"uploader"`
}

type ReleasesResponse struct {
	Url             string    `json:"url"`
	HtmlUrl         string    `json:"html_url"`
	AssetsUrl       string    `json:"assets_url"`
	UploadUrl       string    `json:"upload_url"`
	TarballUrl      string    `json:"tarball_url"`
	ZipballUrl      string    `json:"zipball_url"`
	DiscussionUrl   string    `json:"discussion_url"`
	Id              int       `json:"id"`
	NodeId          string    `json:"node_id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	CreatedAt       time.Time `json:"created_at"`
	PublishedAt     time.Time `json:"published_at"`
	Author          Arthor    `json:"author"`
	Assets          []Asset   `json:"assets"`
}

func (gh *Client) TagAndRelease(branch, tag string) (*ReleasesResponse, error) {
	uri := fmt.Sprintf(epRelease, gh.Host, gh.Org, gh.Repository)
	body := &Release{
		Name:         tag + " - " + time.Now().Format("2006-01-02"),
		TagName:      tag,
		TargetCommit: branch,
	}

	bodyBits, err1 := json.Marshal(body)
	if err1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, err1.Error())
	}

	log.Logf("attempting to publish a release to %v\n", uri)

	res, err2 := gh.Send(uri, "POST", bytes.NewReader(bodyBits))
	if err2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, err2)
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf(
			stderr.ResponseStatusCode,
			res.StatusCode,
			uri,
			res.Status,
		)
	}

	b, err3 := io.ReadAll(res.Body)
	if err3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, err3.Error())
	}

	rr := &ReleasesResponse{}
	err4 := json.Unmarshal(b, rr)
	if err4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, err4.Error())
	}

	return rr, nil
}
