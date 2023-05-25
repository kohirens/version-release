package github

import "time"

type Verification struct {
	Verified  bool        `json:"verified"`
	Reason    string      `json:"reason"`
	Signature interface{} `json:"signature"`
	Payload   interface{} `json:"payload"`
}

type Tree struct {
	Sha string `json:"sha"`
	Url string `json:"url"`
}

type Author struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Commit struct {
	Author       Author       `json:"author"`
	Committer    Author       `json:"committer"`
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
