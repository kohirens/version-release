package git

var stderr = struct {
	CouldNotAddOrigin        string
	CouldNotCheckoutBranch   string
	CouldNotCommit           string
	CouldNotDisplayGitStatus string
	CouldNotGetRemoteUrl     string
	CouldNotGitListRemote    string
	CouldNotPush             string
	CouldNotRemoveOrigin     string
	CouldNotSetGlobalConfig  string
	CouldNotSetRemoteUrl     string
	GitDescribeContains      string
	LastLog                  string
}{
	CouldNotAddOrigin:        "problem adding the origin %s: %s, %s",
	CouldNotCheckoutBranch:   "could not checkout branch: %s; %v",
	CouldNotGetRemoteUrl:     "problem getting the remote push URL: %s, %s",
	CouldNotDisplayGitStatus: "cannot display git status: %s; %v",
	CouldNotGitListRemote:    "could not ls-remote: %s; %v",
	CouldNotCommit:           "could not commit: %s; %s",
	CouldNotPush:             "cannot push changes: %s; %s",
	CouldNotRemoveOrigin:     "problem removing the origin %s: %s, %s",
	CouldNotSetGlobalConfig:  "could not set global config %s; %v",
	CouldNotSetRemoteUrl:     "problem setting the remote push URL: %s, %s",
	GitDescribeContains:      "cannot describe commit %v: %v\n",
	LastLog:                  "could not get latest log: %s",
}

var stdout = struct {
	FoundRemoteBranch  string
	NoTags             string
	SetGitGlobalConfig string
	StagedFiles        string
	Status             string
	TagsInfo           string
}{
	FoundRemoteBranch:  "found remote branch %s\n",
	NoTags:             "no tag for %v was found",
	SetGitGlobalConfig: "set git config global %s\n",
	StagedFiles:        "staged files %s\n",
	Status:             "status: %s",
	TagsInfo:           "tag(s) found %s\n",
}
