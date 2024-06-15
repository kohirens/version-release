package github

var stderr = struct {
	BranchExists             string
	CouldNotGetRequest       string
	CouldNotJsonEncode       string
	CouldNotJsonDecode       string
	CouldNotMakeRequest      string
	CouldNotMergePullRequest string
	CouldNotPingMergeStatus  string
	CouldNotPrepareRequest   string
	CouldNotPostRequest      string
	CouldNotReadResponse     string
	GetMeta                  string
	GetMetaBody              string
	GetMetaUnmarshall        string
	MergeWaitTimeout         string
	ResponseStatusCode       string
}{
	BranchExists:             "the branch %q exists on %s, please delete it manually, then re-run this job so it can complete successfully",
	CouldNotGetRequest:       "could not GET request: %v",
	CouldNotJsonEncode:       "could not encode %t to JSON: %v",
	CouldNotJsonDecode:       "could not decode JSON: %s",
	CouldNotMakeRequest:      "could not make %s request: %s",
	CouldNotMergePullRequest: "unable to merge pr %d: %s",
	CouldNotPingMergeStatus:  "unable to ping pull request %d merge status: %s",
	CouldNotPrepareRequest:   "could not prepare a request: %v",
	CouldNotPostRequest:      "could not POST request: %v",
	CouldNotReadResponse:     "could not read response body: %v",
	GetMeta:                  "could not get github metadata: %v",
	GetMetaBody:              "could not read github meta: %v",
	GetMetaUnmarshall:        "could not decode response from github meta endpoint: %v",
	MergeWaitTimeout:         "pr %d has not merged after %d seconds",
	ResponseStatusCode:       "got a %d response from %v: %s",
}

var stdout = struct {
	CheckMergeStatus  string
	MakePullRequest   string
	MergeResponse     string
	PullRequestMade   string
	PullRequestMerged string
	SendMergeRequest  string
}{
	CheckMergeStatus:  "checking pr %d merged status",
	MakePullRequest:   "making a pull request to %s",
	MergeResponse:     "merge status: %v",
	PullRequestMade:   "pull request %d has been made",
	PullRequestMerged: "pr %d has has been merged",
	SendMergeRequest:  "sending a merge request to %s",
}
