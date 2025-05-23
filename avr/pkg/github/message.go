package github

var stderr = struct {
	BranchExists,
	BuildJWT,
	CouldNotGetRequest,
	CouldNotJsonEncode,
	CouldNotJsonDecode,
	CouldNotMakeRequest,
	CouldNotMergePullRequest,
	CouldNotPingMergeStatus,
	CouldNotPrepareRequest,
	CouldNotPostRequest,
	CouldNotReadResponse,
	EnvNotSet,
	GetMeta,
	GetMetaBody,
	GetMetaUnmarshall,
	GetRef,
	GetTree,
	InstallationToken,
	MakeTree,
	MergeWaitTimeout,
	NewCommit,
	NewRef,
	OpenFile,
	ReadFile,
	ResponseStatusCode,
	UpdateRef string
}{
	BranchExists:             "the branch %q exists on %v, please delete it manually, then re-run this job so it can complete successfully",
	BuildJWT:                 "failed to build a JWT: %v",
	CouldNotGetRequest:       "could not GET request: %v",
	CouldNotJsonEncode:       "could not encode %t to JSON: %v",
	CouldNotJsonDecode:       "could not decode JSON: %s",
	CouldNotMakeRequest:      "could not make %s request: %s",
	CouldNotMergePullRequest: "unable to merge pr %d: %s",
	CouldNotPingMergeStatus:  "unable to ping pull request %d merge status: %s",
	CouldNotPrepareRequest:   "could not prepare a request: %v",
	CouldNotPostRequest:      "could not POST request: %v",
	CouldNotReadResponse:     "could not read response body: %v",
	EnvNotSet:                "environment variable %v is not set",
	GetMeta:                  "could not get github metadata: %v",
	GetMetaBody:              "could not read github meta: %v",
	GetMetaUnmarshall:        "could not decode response from github meta endpoint: %v",
	GetRef:                   "failed to get git reference, responded with status %v %v",
	GetTree:                  "failed to get git tree, responded with status %v %v",
	InstallationToken:        "failed to retrieve an installation token: %v",
	MakeTree:                 "failed to make tree, responded with status %v %v",
	MergeWaitTimeout:         "pr %d has not merged after %d seconds",
	NewCommit:                "failed to make a new commit, responded with status %v %v",
	NewRef:                   "failed to make a new git reference, responded with status %v %v",
	OpenFile:                 "could not open file %v: %v",
	ReadFile:                 "could not read file %v: %v",
	ResponseStatusCode:       "got a %d response from %v: %s",
	UpdateRef:                "failed to update git reference, responded with status %v %v",
}

var stdout = struct {
	BranchExist,
	CheckMergeStatus,
	GetRef,
	GitHubOutputsFile,
	GitOldUrl,
	GitUrl,
	MakePullRequest,
	MergeResponse,
	ParentBranch,
	PullRequestMade,
	PullRequestMerged,
	RemoteBranchStatus,
	RepoUrl,
	RequestUrl,
	SendMergeRequest,
	TriggerWorkflow,
	UpdateRef string
}{
	BranchExist:        "branch exist: %v",
	CheckMergeStatus:   "checking merged status of %v",
	GetRef:             "git reference: %v",
	GitHubOutputsFile:  "begin dumping file %v contents:\n%v\nend contents\n",
	GitOldUrl:          "replacing git remote old url %v",
	GitUrl:             "setting git remote url to %v",
	MakePullRequest:    "making a pull request to %v",
	MergeResponse:      "merge status: %v",
	ParentBranch:       "parent branch %v",
	PullRequestMade:    "pull request %d has been made",
	PullRequestMerged:  "pr %d has has been merged",
	RemoteBranchStatus: "remote branch status:\n%v",
	RepoUrl:            "repo url: %v",
	RequestUrl:         "making a request to: %v",
	SendMergeRequest:   "sending a merge request to %s",
	TriggerWorkflow:    "trigger workflow %v",
	UpdateRef:          "Update reference %v",
}
