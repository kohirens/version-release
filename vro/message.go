package main

var stderr = struct {
	ChgLogAutoUpdate      string
	CommitAlreadyTagged   string
	FiveArgsRequired      string
	InvalidCommit         string
	InvalidSemVer         string
	MissingEnv            string
	NoChangelogChanges    string
	NothingToTag          string
	OpenFile              string
	PublishChangelogArgs  string
	PublishReleaseTagArgs string
}{
	ChgLogAutoUpdate:      "abort, the last commit contains an auto update to the CHANGELOG.md",
	CommitAlreadyTagged:   "commit %v is already tagged",
	FiveArgsRequired:      "5 arguments are required to run this command, see -help",
	InvalidCommit:         "invalid commit given %v",
	InvalidSemVer:         "invalid semantic version given %v",
	MissingEnv:            "%s environment variable is not set",
	NoChangelogChanges:    "the changelog has no changes to be committed",
	NothingToTag:          "next semantic version is empty, so we cannot release a tag",
	OpenFile:              "could not open file %v: %v",
	PublishChangelogArgs:  "3 arguments are required to run this command, see -help",
	PublishReleaseTagArgs: "3 arguments are required to run this command, see -help",
}

var stdout = struct {
	ChgLogUpToDate    string
	ChgLogNotUpToDate string
	CurrentVersion    string
	CurrentVer        string
	GitStatus         string
	Match             string
	NoChanges         string
	ReleaseTag        string
	StartWorkflow     string
	TriggerWorkflow   string
}{
	ChgLogUpToDate:    "the changelog is up to date",
	ChgLogNotUpToDate: "the changelog is not up to date\nchangelog status:\n%s",
	CurrentVersion:    "%v, %v",
	CurrentVer:        "current version %v",
	GitStatus:         "git status output = %s",
	Match:             "entry for %v in the changelog was found, we assume this means the changelog is up-to-date",
	NoChanges:         "no changes to release",
	ReleaseTag:        "releasing %v",
	StartWorkflow:     "starting %v workflow",
	TriggerWorkflow:   "trigger workflow %v",
}

var um = map[string]string{
	"help":                   "display this help",
	"version":                "display version information",
	"semver":                 "set the semantic version for the automated release",
	"tag_and_release_help":   "publish-release-tag -help displays this help",
	"tag_and_release_semver": "provide a semantic version to tag and release",
}
