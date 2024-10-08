package main

var stderr = struct {
	ChgLogAutoUpdate      string
	CommitAlreadyTagged   string
	WorkflowSelectorInput string
	GitHubServer          string
	InvalidCommit         string
	InvalidSemVer         string
	KnownSshKeys          string
	MergeType             string
	MergeTypeEmpty        string
	MissingEnv            string
	NoChangelogChanges    string
	NothingToTag          string
	OpenFile              string
	ParseGitHubRepoEnvVar string
	PublishChangelogArgs  string
	PublishReleaseTagArgs string
	WorkDir               string
}{
	ChgLogAutoUpdate:      "abort, the last commit contains an auto update to the CHANGELOG.md",
	CommitAlreadyTagged:   "commit %v is already tagged",
	WorkflowSelectorInput: "workflow-selector requires arguments changelog file and commit, see -help",
	GitHubServer:          "GitHub server is not set",
	InvalidCommit:         "invalid commit given %v",
	InvalidSemVer:         "invalid semantic version given %v",
	KnownSshKeys:          "KnownSshKeys take the host as the only argument",
	MergeType:             "invalid merge type %v",
	MergeTypeEmpty:        "merge-type was set to an empty string",
	MissingEnv:            "%s environment variable is not set",
	NoChangelogChanges:    "the changelog has no changes to be committed",
	NothingToTag:          "next semantic version is empty, so we cannot release a tag",
	OpenFile:              "could not open file %v: %v",
	ParseGitHubRepoEnvVar: "could not parse %v",
	PublishChangelogArgs:  "path to the changelog is a required argument to run this command, see -help",
	PublishReleaseTagArgs: "3 arguments are required to run this command, see -help",
	WorkDir:               "working directory: %v",
}

var stdout = struct {
	Branch            string
	CciChangelog      string
	ChgLogUpToDate    string
	ChgLogNotUpToDate string
	CurrentVersion    string
	CurrentVer        string
	DbgCommitLog      string
	GaChangelog       string
	GitStatus         string
	Match             string
	NextSemVer        string
	NoChanges         string
	Nothing           string
	ReleaseTag        string
	SemVer            string
	StartWorkflow     string
	Wd                string
}{
	Branch:            "branch %v",
	CciChangelog:      "trigger %v update changelog workflow",
	ChgLogUpToDate:    "the changelog is up to date",
	ChgLogNotUpToDate: "the changelog is not up to date\nchangelog status:\n%s",
	CurrentVersion:    "%v, %v",
	CurrentVer:        "current version %v",
	DbgCommitLog:      "Debug commit log:\n%v",
	GaChangelog:       "GitHub Actions update changelog",
	GitStatus:         "git status output = %s",
	Match:             "entry for %v in the changelog was found, we assume this means the changelog is up-to-date",
	NextSemVer:        "next semVer = %v",
	NoChanges:         "no changes to release",
	Nothing:           "nothing to do, bye!",
	ReleaseTag:        "releasing %v",
	SemVer:            "semantic version set to %v",
	StartWorkflow:     "starting %v workflow",
	Wd:                "working directory %v",
}

var um = map[string]string{
	"help":                            "display this help",
	"version":                         "display version information",
	"cicd":                            "set the CI/CD platform, options are circleci|github (default: circleci).",
	"branch":                          "set the branch to evaluate commit message and or tag",
	"wd":                              "set the working directory of the Git repository to evaluate",
	"gh_api_url":                      "GitHub API URL (default: https://api/github.com)",
	"gh_token":                        "GitHub API access token",
	"semver":                          "set the semantic version for the automated release",
	"publish_changelog_merge_type":    "type of merge to perform when closing the changelog pull request",
	"publish_changelog_github_server": "GitHub server domain name",
}
