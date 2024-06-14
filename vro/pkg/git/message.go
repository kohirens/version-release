package git

var stderr = struct {
	CatFile                  string
	CommitLog                string
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
	WriteCommit              string
}{
	CatFile:                  "git cat-file: %v",
	CommitLog:                "could not get commit log: %s",
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
	GitDescribeContains:      "git describe %s %v",
	LastLog:                  "could not get latest log: %s",
	WriteCommit:              "could not write commit message: %v",
}

var stdout = struct {
	CatFile            string
	FoundRemoteBranch  string
	NoTags             string
	SetGitGlobalConfig string
	StagedFiles        string
	Status             string
	TagsInfo           string
}{
	CatFile:            "git cat-file result: %s",
	FoundRemoteBranch:  "found remote branch %s",
	NoTags:             "no tag for %v was found",
	SetGitGlobalConfig: "set git config global %s",
	StagedFiles:        "staged files %s",
	Status:             "status: %s",
	TagsInfo:           "tag(s) found %s",
}
