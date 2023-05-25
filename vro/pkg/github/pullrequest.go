package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kohirens/stdlib/log"
	"io"
	"time"
)

type Label struct {
	Id          int    `json:"id"`
	NodeId      string `json:"node_id"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

type Team struct {
	Id                  int    `json:"id"`
	NodeId              string `json:"node_id"`
	Url                 string `json:"url"`
	HtmlUrl             string `json:"html_url"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Description         string `json:"description"`
	Privacy             string `json:"privacy"`
	NotificationSetting string `json:"notification_setting"`
	Permission          string `json:"permission"`
	MembersUrl          string `json:"members_url"`
	RepositoriesUrl     string `json:"repositories_url"`
}

type Permissions struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	SpdxId string `json:"spdx_id"`
	NodeId string `json:"node_id"`
}

type Repo struct {
	Id               int         `json:"id"`
	NodeId           string      `json:"node_id"`
	Name             string      `json:"name"`
	FullName         string      `json:"full_name"`
	Owner            Profile     `json:"owner"`
	Private          bool        `json:"private"`
	HtmlUrl          string      `json:"html_url"`
	Description      string      `json:"description"`
	Fork             bool        `json:"fork"`
	Url              string      `json:"url"`
	ArchiveUrl       string      `json:"archive_url"`
	AssigneesUrl     string      `json:"assignees_url"`
	BlobsUrl         string      `json:"blobs_url"`
	BranchesUrl      string      `json:"branches_url"`
	CollaboratorsUrl string      `json:"collaborators_url"`
	CommentsUrl      string      `json:"comments_url"`
	CommitsUrl       string      `json:"commits_url"`
	CompareUrl       string      `json:"compare_url"`
	ContentsUrl      string      `json:"contents_url"`
	ContributorsUrl  string      `json:"contributors_url"`
	DeploymentsUrl   string      `json:"deployments_url"`
	DownloadsUrl     string      `json:"downloads_url"`
	EventsUrl        string      `json:"events_url"`
	ForksUrl         string      `json:"forks_url"`
	GitCommitsUrl    string      `json:"git_commits_url"`
	GitRefsUrl       string      `json:"git_refs_url"`
	GitTagsUrl       string      `json:"git_tags_url"`
	GitUrl           string      `json:"git_url"`
	IssueCommentUrl  string      `json:"issue_comment_url"`
	IssueEventsUrl   string      `json:"issue_events_url"`
	IssuesUrl        string      `json:"issues_url"`
	KeysUrl          string      `json:"keys_url"`
	LabelsUrl        string      `json:"labels_url"`
	LanguagesUrl     string      `json:"languages_url"`
	MergesUrl        string      `json:"merges_url"`
	MilestonesUrl    string      `json:"milestones_url"`
	NotificationsUrl string      `json:"notifications_url"`
	PullsUrl         string      `json:"pulls_url"`
	ReleasesUrl      string      `json:"releases_url"`
	SshUrl           string      `json:"ssh_url"`
	StargazersUrl    string      `json:"stargazers_url"`
	StatusesUrl      string      `json:"statuses_url"`
	SubscribersUrl   string      `json:"subscribers_url"`
	SubscriptionUrl  string      `json:"subscription_url"`
	TagsUrl          string      `json:"tags_url"`
	TeamsUrl         string      `json:"teams_url"`
	TreesUrl         string      `json:"trees_url"`
	CloneUrl         string      `json:"clone_url"`
	MirrorUrl        string      `json:"mirror_url"`
	HooksUrl         string      `json:"hooks_url"`
	SvnUrl           string      `json:"svn_url"`
	Homepage         string      `json:"homepage"`
	Language         interface{} `json:"language"`
	ForksCount       int         `json:"forks_count"`
	StargazersCount  int         `json:"stargazers_count"`
	WatchersCount    int         `json:"watchers_count"`
	Size             int         `json:"size"`
	DefaultBranch    string      `json:"default_branch"`
	OpenIssuesCount  int         `json:"open_issues_count"`
	Topics           []string    `json:"topics"`
	HasIssues        bool        `json:"has_issues"`
	HasProjects      bool        `json:"has_projects"`
	HasWiki          bool        `json:"has_wiki"`
	HasPages         bool        `json:"has_pages"`
	HasDownloads     bool        `json:"has_downloads"`
	HasDiscussions   bool        `json:"has_discussions"`
	Archived         bool        `json:"archived"`
	Disabled         bool        `json:"disabled"`
	PushedAt         time.Time   `json:"pushed_at"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Permissions      Permissions `json:"permissions"`
	AllowRebaseMerge bool        `json:"allow_rebase_merge"`
	TempCloneToken   string      `json:"temp_clone_token"`
	AllowSquashMerge bool        `json:"allow_squash_merge"`
	AllowMergeCommit bool        `json:"allow_merge_commit"`
	AllowForking     bool        `json:"allow_forking"`
	Forks            int         `json:"forks"`
	OpenIssues       int         `json:"open_issues"`
	License          License     `json:"license"`
	Watchers         int         `json:"watchers"`
}

type Milestone struct {
	Url          string    `json:"url"`
	HtmlUrl      string    `json:"html_url"`
	LabelsUrl    string    `json:"labels_url"`
	Id           int       `json:"id"`
	NodeId       string    `json:"node_id"`
	Number       int       `json:"number"`
	State        string    `json:"state"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Creator      Profile   `json:"creator"`
	OpenIssues   int       `json:"open_issues"`
	ClosedIssues int       `json:"closed_issues"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ClosedAt     time.Time `json:"closed_at"`
	DueOn        time.Time `json:"due_on"`
}

type Ref struct {
	Label string  `json:"label"`
	Ref   string  `json:"ref"`
	Sha   string  `json:"sha"`
	User  Profile `json:"user"`
	Repo  Repo    `json:"repo"`
}

type Anchor struct {
	Href string `json:"href"`
}

type Links struct {
	Self           Anchor `json:"self"`
	Html           Anchor `json:"html"`
	Issue          Anchor `json:"issue"`
	Comments       Anchor `json:"comments"`
	ReviewComments Anchor `json:"review_comments"`
	ReviewComment  Anchor `json:"review_comment"`
	Commits        Anchor `json:"commits"`
	Statuses       Anchor `json:"statuses"`
}

type PullRequest struct {
	Url                 string      `json:"url"`
	Id                  int         `json:"id"`
	NodeId              string      `json:"node_id"`
	HtmlUrl             string      `json:"html_url"`
	DiffUrl             string      `json:"diff_url"`
	PatchUrl            string      `json:"patch_url"`
	IssueUrl            string      `json:"issue_url"`
	CommitsUrl          string      `json:"commits_url"`
	ReviewCommentsUrl   string      `json:"review_comments_url"`
	ReviewCommentUrl    string      `json:"review_comment_url"`
	CommentsUrl         string      `json:"comments_url"`
	StatusesUrl         string      `json:"statuses_url"`
	Number              int         `json:"number"`
	State               string      `json:"state"`
	Locked              bool        `json:"locked"`
	Title               string      `json:"title"`
	User                Profile     `json:"user"`
	Body                string      `json:"body"`
	Labels              []Label     `json:"labels"`
	Milestone           Milestone   `json:"milestone"`
	ActiveLockReason    string      `json:"active_lock_reason"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	ClosedAt            time.Time   `json:"closed_at"`
	MergedAt            time.Time   `json:"merged_at"`
	MergeCommitSha      string      `json:"merge_commit_sha"`
	Assignee            Profile     `json:"assignee"`
	Assignees           []Profile   `json:"assignees"`
	RequestedReviewers  []Profile   `json:"requested_reviewers"`
	RequestedTeams      []Team      `json:"requested_teams"`
	Head                Ref         `json:"head"`
	Base                Ref         `json:"base"`
	Links               Links       `json:"_links"`
	AuthorAssociation   string      `json:"author_association"`
	AutoMerge           interface{} `json:"auto_merge"`
	Draft               bool        `json:"draft"`
	Merged              bool        `json:"merged"`
	Mergeable           bool        `json:"mergeable"`
	Rebaseable          bool        `json:"rebaseable"`
	MergeableState      string      `json:"mergeable_state"`
	MergedBy            Profile     `json:"merged_by"`
	Comments            int         `json:"comments"`
	ReviewComments      int         `json:"review_comments"`
	MaintainerCanModify bool        `json:"maintainer_can_modify"`
	Commits             int         `json:"commits"`
	Additions           int         `json:"additions"`
	Deletions           int         `json:"deletions"`
	ChangedFiles        int         `json:"changed_files"`
}

type PullRequestBody struct {
	Base  string `json:"base"`
	Body  string `json:"body,omitempty"`
	Title string `json:"title"`
	Head  string `json:"head"`
}

type PullRequestMergeBody struct {
	CommitTitle   string `json:"commit_title,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	MergeMethod   string `json:"merge_method"`
	PullNumber    int    `json:"pull_number"`
}
type PullRequestMergeResponse struct {
	Sha     string `json:"sha"`
	Merged  bool   `json:"merged"`
	Message string `json:"message"`
}

// MergePullRequest Merge a pull request.
func (gh *Client) MergePullRequest(prNumber int, mergeMethod string) (*PullRequestMergeResponse, error) {
	uri := fmt.Sprintf(epPullMerge, gh.Host, gh.Org, gh.Repository, prNumber)

	body := &PullRequestMergeBody{
		PullNumber:  prNumber,
		MergeMethod: mergeMethod,
	}

	bodyBits, err1 := json.Marshal(body)
	if err1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, err1.Error())
	}

	log.Logf(stdout.SendMergeRequest, uri)

	res, err2 := gh.Send(uri, "PUT", bytes.NewReader(bodyBits))
	if err2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotMakeRequest, "PUT", err2.Error())
	}

	b, err3 := io.ReadAll(res.Body)
	if err3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, err3.Error())
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf(stderr.CouldNotMergePullRequest, prNumber, res.StatusCode)
	}

	pr := &PullRequestMergeResponse{}
	err4 := json.Unmarshal(b, pr)
	if err4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, pr, err4.Error())
	}

	return pr, nil
}

// OpenPullRequest Opens a pull request against the base branch.
// gh pr create --base "${PARAM_MAIN_TRUNK_BRANCH}" --head "${GEN_BRANCH_NAME}" --fill
func (gh *Client) OpenPullRequest(base, branch, title, summary string) (*PullRequest, error) {
	uri := fmt.Sprintf(epPulls, gh.Host, gh.Org, gh.Repository)
	body := &PullRequestBody{
		Base:  base,
		Body:  summary,
		Head:  branch,
		Title: title,
	}

	bodyBits, err1 := json.Marshal(body)
	if err1 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonEncode, body, err1.Error())
	}

	log.Logf(stdout.MakePullRequest, uri)

	res, err2 := gh.Send(uri, "POST", bytes.NewReader(bodyBits))
	if err2 != nil {
		return nil, fmt.Errorf(stderr.CouldNotPostRequest, err2.Error())
	}

	if res.StatusCode != 201 {
		return nil, fmt.Errorf(stderr.ResponseStatusCode, res.StatusCode)
	}

	b, err3 := io.ReadAll(res.Body)
	if err3 != nil {
		return nil, fmt.Errorf(stderr.CouldNotReadResponse, err3.Error())
	}

	pr := &PullRequest{}
	err4 := json.Unmarshal(b, pr)
	if err4 != nil {
		return nil, fmt.Errorf(stderr.CouldNotJsonDecode, err4.Error())
	}

	log.Logf(stdout.PullRequestMade, pr.Number)

	return pr, nil
}
