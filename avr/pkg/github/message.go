package github

var stderr = struct {
	BranchExists,
	CouldNotGetRequest,
	CouldNotJsonEncode,
	CouldNotJsonDecode,
	CouldNotMakeRequest,
	CouldNotMergePullRequest,
	CouldNotPingMergeStatus,
	CouldNotPrepareRequest,
	CouldNotPostRequest,
	CouldNotReadResponse,
	GetMeta,
	GetMetaBody,
	GetMetaUnmarshall,
	MergeWaitTimeout,
	OpenFile,
	ReadFile,
	ResponseStatusCode string
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
	OpenFile:                 "could not open file %v: %v",
	ReadFile:                 "could not read file %v: %v",
	ResponseStatusCode:       "got a %d response from %v: %s",
}

var stdout = struct {
	BranchExist,
	CheckMergeStatus,
	GitHubOutputsFile,
	GitOldUrl,
	GitUrl,
	MakePullRequest,
	MergeResponse,
	PullRequestMade,
	PullRequestMerged,
	RemoteBranchStatus,
	RepoUrl,
	RequestUrl,
	SendMergeRequest,
	TriggerWorkflow string
}{
	BranchExist:        "branch exist: %v",
	CheckMergeStatus:   "checking merged status of %v",
	GitHubOutputsFile:  "begin dumping file %v contents:\n%v\nend contents\n",
	GitOldUrl:          "replacing git remote old url %v",
	GitUrl:             "setting git remote url to %v",
	MakePullRequest:    "making a pull request to %v",
	MergeResponse:      "merge status: %v",
	PullRequestMade:    "pull request %d has been made",
	PullRequestMerged:  "pr %d has has been merged",
	RemoteBranchStatus: "remote branch status:\n%v",
	RepoUrl:            "repo url: %v",
	RequestUrl:         "making a request to: %v",
	SendMergeRequest:   "sending a merge request to %s",
	TriggerWorkflow:    "trigger workflow %v",
}
